/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"sort"

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
				return executeUpdate(repo, args[0], true)
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

// undoCmd returns the command to mark a task as not completed.
func undoCmd(repo todo.TodoStore) *cobra.Command {
	return &cobra.Command{
		Use:     "undo",
		Aliases: []string{"undone"},
		Short:   "Mark a task as not completed.",
		Long: `Mark a done task as not done, using its number from 'tri list'.
Example: tri undo 1`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				slog.Error("No index provided")
				return
			}

			err := WaitSpinner("Marking as not done", func() error {
				return executeUpdate(repo, args[0], false)
			})

			if err != nil {
				fmt.Println("Error: Failed to update the task. Please ensure you provided a valid index.")
				return
			}
		},
	}
}

// executeUpdate handles the logic for finding a task and setting its done state.
func executeUpdate(r todo.TodoStore, argStr string, done bool) error {
	items, err := r.ListItems()
	if err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}

	// Match the ordering the user saw in `tri list`.
	sort.Sort(todo.ByPri(items))

	idx, err := validateIndex(argStr, len(items))
	if err != nil {
		return err
	}

	err = r.UpdateItemStatus(items[idx-1].ID, done)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	slog.Info("Marked_done_success", "index", idx)
	return nil
}
