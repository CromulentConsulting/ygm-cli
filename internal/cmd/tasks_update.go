package cmd

import (
	"fmt"
	"strconv"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	updateTitle       string
	updateDescription string
	updateStatus      string
)

var tasksUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a marketing task",
	Long: `Update title, description, or status of an existing marketing task.

Examples:
  ygm tasks update 42 --title "New title"
  ygm tasks update 42 --status completed
  ygm tasks update 42 --title "Updated" --description "New description" --status in_progress`,
	Args: cobra.ExactArgs(1),
	RunE: runTasksUpdate,
}

func init() {
	tasksUpdateCmd.Flags().StringVar(&updateTitle, "title", "", "New task title")
	tasksUpdateCmd.Flags().StringVar(&updateDescription, "description", "", "New task description")
	tasksUpdateCmd.Flags().StringVar(&updateStatus, "status", "", "New status (pending, in_progress, completed)")
}

func runTasksUpdate(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	if updateTitle == "" && updateDescription == "" && updateStatus == "" {
		return fmt.Errorf("at least one of --title, --description, or --status is required")
	}

	account, err := getActiveAccount()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.APIURL, account.Token)

	req := api.UpdateTaskRequest{
		Title:       updateTitle,
		Description: updateDescription,
		Status:      updateStatus,
	}

	task, err := client.UpdateTask(id, req)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if jsonOutput {
		return outputJSON(task)
	}

	fmt.Printf("Updated task #%d: %s\n", task.ID, task.Title)

	return nil
}
