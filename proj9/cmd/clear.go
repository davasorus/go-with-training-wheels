/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// clearCmd returns the command to purge all completed tasks.
func clearCmd(repo todo.TodoStore) *cobra.Command {
	return &cobra.Command{
		Use:     "clear",
		Aliases: []string{"clr"},
		Short:   "Purge all items marked as 'Done'.",
		Long:    `Remove all tasks from the database that are already marked as completed.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := WaitSpinner("Clearing completed tasks", func() error {
				return executeClear(repo)
			})

			if err != nil {
				fmt.Println("Error: Failed to clear completed items.")
				return
			}
			fmt.Println("Successfully cleared all completed items")
		},
	}
}

// executeClear handles the logic for removing all finished tasks.
func executeClear(r todo.TodoStore) error {
	items, err := r.ListItems()
	if err != nil {
		return fmt.Errorf("failed to fetch items: %w", err)
	}

	for _, item := range items {
		if item.Done {
			err := r.DeleteItem(item.ID)
			if err != nil {
				slog.Error("Failed to delete completed item", "id", item.ID, "error", err)
				// Continue with other items even if one fails
			}
		}
	}

	return nil
}
