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

// doneCmd returns the command to mark a task as completed.
func doneCmd(repo todo.TodoStore) *cobra.Command {
	return &cobra.Command{
		Use:     "done",
		Aliases: []string{"do"},
		Short:   "Mark a task as completed.",
		Long: `Mark a task as done using its position in the list.
Example: tri do 1`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				slog.Error("No index provided")
				return
			}

			err := WaitSpinner("Marking as done", func() error {
				return executeUpdate(repo, args[0])
			})

			if err != nil {
				fmt.Println("Error: Failed to mark the task as done. Please ensure you provided a valid index.")
				return
			}
		},
	}
}

func init() {
	// Note: doneCmd is now added in root.go by passing the repository.
}

// executeUpdate handles the logic for finding and marking a task as completed.
func executeUpdate(r todo.TodoStore, argStr string) error {
	items, err := r.ListItems()
	if err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}

	idx, err := validateIndex(argStr, len(items))
	if err != nil {
		return err
	}

	err = r.UpdateItemStatus(idx-1, true)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	slog.Info("Marked_done_success", "index", idx)
	return nil
}
