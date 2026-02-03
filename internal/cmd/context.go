package cmd

import (
	"fmt"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Get full context for AI prompts",
	Long: `Get a full context dump suitable for AI coding assistants.

This includes your brand DNA, marketing plan, and pending tasks
in a JSON format that can be included in prompts for AI tools
like Claude, ChatGPT, or GitHub Copilot.`,
	RunE: runContext,
}

func runContext(cmd *cobra.Command, args []string) error {
	account, err := getActiveAccount()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.APIURL, account.Token)
	ctx, err := client.GetContext()
	if err != nil {
		return fmt.Errorf("failed to fetch context: %w", err)
	}

	// Context always outputs JSON (it's designed for machine consumption)
	return outputJSON(ctx)
}
