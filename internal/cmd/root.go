package cmd

import (
	"fmt"
	"os"

	"github.com/CromulentConsulting/ygm-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Flags
	orgFlag    string
	jsonOutput bool

	// Global config
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "ygm",
	Short: "You've Got Marketing CLI",
	Long: `ygm is a command-line interface for You've Got Marketing.

It provides access to your brand DNA, marketing tasks, and context
for use with AI coding assistants.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for login command
		if cmd.Name() == "login" || cmd.Name() == "version" {
			return nil
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg == nil || len(cfg.Accounts) == 0 {
			fmt.Fprintln(os.Stderr, "Not logged in. Run 'ygm login' first.")
			os.Exit(1)
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&orgFlag, "org", "", "Organization to use (overrides default)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(brandCmd)
	rootCmd.AddCommand(tasksCmd)
	rootCmd.AddCommand(contextCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ygm version %s\n", Version)
	},
}

// getActiveAccount returns the account to use based on --org flag or default
func getActiveAccount() (*config.Account, error) {
	if cfg == nil {
		return nil, fmt.Errorf("not logged in")
	}

	orgSlug := orgFlag
	if orgSlug == "" {
		orgSlug = cfg.DefaultOrg
	}

	if orgSlug == "" {
		// Return first account if no default set
		for _, account := range cfg.Accounts {
			return &account, nil
		}
		return nil, fmt.Errorf("no accounts configured")
	}

	account, ok := cfg.Accounts[orgSlug]
	if !ok {
		return nil, fmt.Errorf("organization '%s' not found in config", orgSlug)
	}

	return &account, nil
}
