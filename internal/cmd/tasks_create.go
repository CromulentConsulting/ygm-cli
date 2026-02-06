package cmd

import (
	"fmt"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	taskTitle       string
	taskDescription string
	taskPlatform    string
	taskAssetType   string
	taskDate        string
)

var tasksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new marketing task",
	Long: `Create a new marketing task in your current marketing plan.

Examples:
  ygm tasks create --title "Post on Reddit" --platform reddit
  ygm tasks create --title "Launch tweet" --description "Announce v2" --platform twitter --date 2026-02-11
  ygm tasks create --title "Blog post" --platform blog --asset-type copy --json`,
	RunE: runTasksCreate,
}

func init() {
	tasksCreateCmd.Flags().StringVar(&taskTitle, "title", "", "Task title (required)")
	tasksCreateCmd.Flags().StringVar(&taskDescription, "description", "", "Task description")
	tasksCreateCmd.Flags().StringVar(&taskPlatform, "platform", "", "Platform (twitter, instagram, linkedin, reddit, etc.)")
	tasksCreateCmd.Flags().StringVar(&taskAssetType, "asset-type", "", "Asset type (image, copy, video)")
	tasksCreateCmd.Flags().StringVar(&taskDate, "date", "", "Suggested post date (YYYY-MM-DD)")
	tasksCreateCmd.MarkFlagRequired("title")
}

func runTasksCreate(cmd *cobra.Command, args []string) error {
	account, err := getActiveAccount()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.APIURL, account.Token)

	req := api.CreateTaskRequest{
		Title:       taskTitle,
		Description: taskDescription,
		Platform:    taskPlatform,
		AssetType:   taskAssetType,
	}
	if taskDate != "" {
		req.SuggestedPostDate = &taskDate
	}

	task, err := client.CreateTask(req)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	if jsonOutput {
		return outputJSON(task)
	}

	fmt.Printf("Created task #%d: %s\n", task.ID, task.Title)
	if task.Platform != "" {
		fmt.Printf("  Platform: %s\n", task.Platform)
	}
	if task.SuggestedPostDate != nil {
		fmt.Printf("  Date: %s\n", *task.SuggestedPostDate)
	}

	return nil
}
