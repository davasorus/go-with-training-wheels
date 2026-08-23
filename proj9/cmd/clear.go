/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// clearCmd returns the command to purge all completed tasks.
func clearCmd(repo todo.TodoStore) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:     "clear",
		Aliases: []string{"clr"},
		Short:   "Purge all items marked as 'Done'.",
		Long:    `Remove all tasks from the database that are already marked as completed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := repo.ListItems()
			if err != nil {
				return fmt.Errorf("failed to load items: %w", err)
			}
			doneCount := 0
			for _, item := range items {
				if item.Done {
					doneCount++
				}
			}
			if doneCount == 0 {
				fmt.Println("Nothing to clear: no completed tasks.")
				return nil
			}

			if !yes && !confirm(fmt.Sprintf("Delete %d completed task(s)?", doneCount)) {
				fmt.Println("Canceled.")
				return nil
			}

			err = WaitSpinner("Clearing completed tasks", func() error {
				return executeClear(repo)
			})
			if err != nil {
				return err
			}
			fmt.Printf("Cleared %d completed task(s).\n", doneCount)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Clear without asking for confirmation.")
	return cmd
}

// executeClear handles the logic for removing all finished tasks.
func executeClear(r todo.TodoStore) error {
	items, err := r.ListItems()
	if err != nil {
		return fmt.Errorf("failed to fetch items: %w", err)
	}

	for _, item := range items {
		if item.Done {
			if err := r.DeleteItem(item.ID); err != nil {
				return fmt.Errorf("failed to delete completed item %d: %w", item.ID, err)
			}
		}
	}

	return nil
}
