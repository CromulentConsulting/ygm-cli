package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
	"github.com/CromulentConsulting/ygm-cli/internal/auth"
	"github.com/CromulentConsulting/ygm-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	apiURLFlag   string
	tokenNameFlag string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with You've Got Marketing",
	Long: `Start the device flow authentication to connect the CLI to your account.

This will open your browser where you can enter a code and authorize the CLI.
Once authorized, the token will be saved locally for future use.`,
	RunE: runLogin,
}

func init() {
	loginCmd.Flags().StringVar(&apiURLFlag, "api-url", config.DefaultAPIURL, "API URL (for development)")
	loginCmd.Flags().StringVar(&tokenNameFlag, "name", "", "Name for this token (e.g., 'MacBook CLI')")
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Load or create config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg == nil {
		cfg = config.NewConfig()
	}

	// Override API URL if specified
	if apiURLFlag != "" {
		cfg.APIURL = apiURLFlag
	}

	fmt.Println("Starting device flow authentication...")
	fmt.Println()

	// Request device code
	deviceFlow := auth.NewDeviceFlow(cfg.APIURL)
	deviceCode, err := deviceFlow.RequestDeviceCode()
	if err != nil {
		return fmt.Errorf("failed to start authentication: %w", err)
	}

	// Display the user code
	fmt.Printf("Your code: %s\n", deviceCode.UserCode)
	fmt.Println()
	fmt.Printf("Opening %s in your browser...\n", deviceCode.VerificationURI)
	fmt.Println()

	// Try to open browser
	codeURL := fmt.Sprintf("%s?code=%s", deviceCode.VerificationURI, deviceCode.UserCode)
	if err := auth.OpenBrowser(codeURL); err != nil {
		fmt.Fprintf(os.Stderr, "Could not open browser automatically.\n")
		fmt.Printf("Please visit: %s\n", codeURL)
		fmt.Println()
	}

	fmt.Println("Waiting for authorization...")

	// Set up timeout
	timeout := time.After(time.Duration(deviceCode.ExpiresIn) * time.Second)

	type pollResult struct {
		token *api.TokenResponse
		err   error
	}
	done := make(chan pollResult)

	// Poll in background
	go func() {
		// Generate token name
		name := tokenNameFlag
		if name == "" {
			hostname, _ := os.Hostname()
			if hostname == "" {
				hostname = "CLI"
			}
			name = fmt.Sprintf("%s %s", hostname, time.Now().Format("2006-01-02"))
		}

		token, err := deviceFlow.PollForToken(deviceCode.DeviceCode, deviceCode.Interval, name)
		done <- pollResult{token, err}
	}()

	// Wait for result or timeout
	select {
	case result := <-done:
		if result.err != nil {
			return fmt.Errorf("authentication failed: %w", result.err)
		}

		token := result.token

		// Save to config
		cfg.AddAccount(token.Organization.Slug, config.Account{
			Token:     token.AccessToken,
			UserEmail: token.User.Email,
			OrgID:     token.Organization.ID,
			OrgName:   token.Organization.Name,
		})

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println()
		fmt.Println("Successfully authenticated!")
		fmt.Printf("  Organization: %s\n", token.Organization.Name)
		fmt.Printf("  User: %s\n", token.User.Email)
		fmt.Println()
		fmt.Println("You can now use 'ygm brand', 'ygm tasks', and 'ygm context'.")

		return nil

	case <-timeout:
		return fmt.Errorf("authentication timed out - please try again")
	}
}
