package cmd

import (
	"fmt"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	statusFilter   string
	platformFilter string
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List marketing tasks",
	Long: `List marketing tasks from your marketing plan.

Tasks include content creation items like social media posts,
blog articles, and other marketing materials.`,
	RunE: runTasks,
}

func init() {
	tasksCmd.Flags().StringVar(&statusFilter, "status", "", "Filter by status (pending, in_progress, completed)")
	tasksCmd.Flags().StringVar(&platformFilter, "platform", "", "Filter by platform (instagram, twitter, linkedin, etc.)")
	tasksCmd.AddCommand(tasksCreateCmd)
}

func runTasks(cmd *cobra.Command, args []string) error {
	account, err := getActiveAccount()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.APIURL, account.Token)
	tasks, err := client.GetTasks(statusFilter, platformFilter)
	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	if jsonOutput {
		return outputJSON(map[string]interface{}{"tasks": tasks})
	}

	return outputTasksText(tasks)
}

func outputTasksText(tasks []api.Task) error {
	fmt.Printf("Tasks (%d total)\n", len(tasks))
	fmt.Println("================")
	fmt.Println()

	// Group by status
	pending := []api.Task{}
	inProgress := []api.Task{}
	completed := []api.Task{}

	for _, t := range tasks {
		switch t.Status {
		case "pending":
			pending = append(pending, t)
		case "in_progress":
			inProgress = append(inProgress, t)
		case "completed", "shared":
			completed = append(completed, t)
		}
	}

	if len(inProgress) > 0 {
		fmt.Println("In Progress:")
		for _, t := range inProgress {
			printTask(t)
		}
		fmt.Println()
	}

	if len(pending) > 0 {
		fmt.Println("Pending:")
		for _, t := range pending {
			printTask(t)
		}
		fmt.Println()
	}

	if len(completed) > 0 {
		fmt.Println("Completed:")
		for _, t := range completed {
			printTask(t)
		}
		fmt.Println()
	}

	return nil
}

func printTask(t api.Task) {
	platform := t.Platform
	if platform == "" {
		platform = "general"
	}
	fmt.Printf("  [%d] %s\n", t.ID, t.Title)
	fmt.Printf("      Platform: %s", platform)
	if t.SuggestedPostDate != nil {
		fmt.Printf(" | Date: %s", *t.SuggestedPostDate)
	}
	fmt.Println()
	if t.Description != "" {
		// Truncate description if too long
		desc := t.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Printf("      %s\n", desc)
	}
}
