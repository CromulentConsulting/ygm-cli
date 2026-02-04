package cmd

import (
	"fmt"
	"os"

	"github.com/CromulentConsulting/ygm-cli/internal/config"
	"github.com/CromulentConsulting/ygm-cli/internal/skills"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link [org-slug]",
	Short: "Link this project to an organization",
	Long: `Link the current directory to a specific organization.

This creates a .ygm.yml file in the current directory that specifies
which organization to use for this project. This is useful when you
work with multiple organizations/clients.

The .ygm.yml file can be committed to version control (it contains
no secrets, just the org identifier).

Precedence for organization selection:
  1. --org flag (highest)
  2. .ygm.yml (project-specific)
  3. default_org in ~/.config/ygm/config.yml`,
	Example: `  # Link to an org interactively (shows available orgs)
  ygm link

  # Link to a specific org
  ygm link acme-corp`,
	RunE: runLink,
}

func runLink(cmd *cobra.Command, args []string) error {
	// Load global config to get available accounts
	globalCfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if globalCfg == nil || len(globalCfg.Accounts) == 0 {
		fmt.Fprintln(os.Stderr, "Not logged in. Run 'ygm login' first.")
		os.Exit(1)
	}

	var orgSlug string

	if len(args) > 0 {
		// Org specified as argument
		orgSlug = args[0]

		// Validate it exists
		if _, ok := globalCfg.Accounts[orgSlug]; !ok {
			fmt.Fprintf(os.Stderr, "Organization '%s' not found.\n\n", orgSlug)
			fmt.Fprintln(os.Stderr, "Available organizations:")
			for slug, account := range globalCfg.Accounts {
				fmt.Fprintf(os.Stderr, "  - %s (%s)\n", slug, account.OrgName)
			}
			os.Exit(1)
		}
	} else {
		// Interactive: show available orgs
		if len(globalCfg.Accounts) == 1 {
			// Only one account, use it
			for slug := range globalCfg.Accounts {
				orgSlug = slug
			}
		} else {
			fmt.Println("Available organizations:")
			fmt.Println()
			i := 1
			slugs := make([]string, 0, len(globalCfg.Accounts))
			for slug, account := range globalCfg.Accounts {
				fmt.Printf("  %d. %s (%s)\n", i, slug, account.OrgName)
				slugs = append(slugs, slug)
				i++
			}
			fmt.Println()
			fmt.Print("Enter org slug to link: ")

			var input string
			fmt.Scanln(&input)

			if input == "" {
				return fmt.Errorf("no organization selected")
			}

			// Check if valid
			if _, ok := globalCfg.Accounts[input]; !ok {
				return fmt.Errorf("organization '%s' not found", input)
			}
			orgSlug = input
		}
	}

	// Check if already linked
	existingPath, _ := config.LocalConfigPath()
	if existingPath != "" {
		existingCfg, _ := config.LoadLocal()
		if existingCfg != nil && existingCfg.Org == orgSlug {
			fmt.Printf("Already linked to '%s'\n", orgSlug)
			return nil
		}
		if existingCfg != nil {
			fmt.Printf("Updating link from '%s' to '%s'\n", existingCfg.Org, orgSlug)
		}
	}

	// Create local config
	localCfg := &config.LocalConfig{
		Org: orgSlug,
	}

	if err := localCfg.Save(); err != nil {
		return fmt.Errorf("failed to save local config: %w", err)
	}

	account := globalCfg.Accounts[orgSlug]
	fmt.Printf("Linked to '%s' (%s)\n", orgSlug, account.OrgName)
	fmt.Printf("Created %s\n", config.LocalConfigFile)

	// Install local agent skills for AI assistant discovery
	if err := skills.InstallLocal(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not install local agent skills: %v\n", err)
	}

	return nil
}
