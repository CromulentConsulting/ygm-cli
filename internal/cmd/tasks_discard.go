package cmd

import (
	"fmt"
	"strconv"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
	"github.com/spf13/cobra"
)

var tasksDiscardCmd = &cobra.Command{
	Use:   "discard <id>",
	Short: "Discard a marketing task",
	Long: `Soft-delete a marketing task. The task can be restored later if needed.

Examples:
  ygm tasks discard 42
  ygm tasks discard 42 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksDiscard,
}

func runTasksDiscard(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	account, err := getActiveAccount()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.APIURL, account.Token)

	result, err := client.DiscardTask(id)
	if err != nil {
		return fmt.Errorf("failed to discard task: %w", err)
	}

	if jsonOutput {
		return outputJSON(result)
	}

	fmt.Printf("Discarded task #%d\n", id)

	return nil
}
