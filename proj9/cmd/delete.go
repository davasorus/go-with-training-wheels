/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// validateIndex validates the index provided and returns it as an int.
func validateIndex(input string, count int) (int, error) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("not a number")
	}

	if i <= 0 || i > count {
		return 0, fmt.Errorf("%d does not match any item", i)
	}

	return i, nil
}

// deleteCmd returns the command to remove a task.
func deleteCmd(repo todo.TodoStore) *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "Remove a task from the list.",
		Long: `Remove a task from the list based on its position or ID.
Example: tri delete 1`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Error: No input provided. Please specify an index (e.g., 'tri delete 1').")
				return
			}

			err := WaitSpinner("Deleting item", func() error {
				return executeDelete(repo, args[0])
			})

			if err != nil {
				fmt.Println("Error: Failed to delete the item.")
				return
			}
		},
	}
}

// executeDelete handles the logic for finding and removing a task.
func executeDelete(r todo.TodoStore, argStr string) error {
	val, err := strconv.Atoi(argStr)
	if err != nil {
		return fmt.Errorf("invalid input: %s is not a number", argStr)
	}

	items, err := r.ListItems()
	if err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}

	// Match the ordering the user saw in `tri list`.
	sort.Sort(todo.ByPri(items))

	var idToDelete int
	if val >= 1 && val <= len(items) {
		idToDelete = items[val-1].ID
	} else {
		// Fall back to treating the value as a raw database ID.
		idToDelete = val
	}

	err = r.DeleteItem(idToDelete)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	slog.Info("Removed_task_successfully", "input", argStr, "id", idToDelete)
	return nil
}
