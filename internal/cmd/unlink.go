package cmd

import (
	"fmt"

	"github.com/CromulentConsulting/ygm-cli/internal/config"
	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink this project from its organization",
	Long: `Remove the .ygm.yml file from the current project.

After unlinking, the CLI will use the default organization from your
global config (~/.config/ygm/config.yml) or require --org flag.`,
	RunE: runUnlink,
}

func runUnlink(cmd *cobra.Command, args []string) error {
	path, err := config.LocalConfigPath()
	if err != nil {
		return err
	}

	if path == "" {
		fmt.Println("No local config found. Project is not linked.")
		return nil
	}

	// Show what we're removing
	localCfg, _ := config.LoadLocal()
	if localCfg != nil {
		fmt.Printf("Unlinking from '%s'\n", localCfg.Org)
	}

	if err := config.RemoveLocal(); err != nil {
		return fmt.Errorf("failed to remove local config: %w", err)
	}

	fmt.Printf("Removed %s\n", path)
	return nil
}
