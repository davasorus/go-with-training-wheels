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

// doneCmd returns the command to mark a task as complete.
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

			err := executeUpdate(repo, args[0])
			if err != nil {
				slog.Error("Failed to mark done", "error", err)
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

	i, err := validateIndex(argStr, len(items))
	if err != nil {
		return err
	}

	items[i-1].Done = true
	sort.Sort(todo.ByPri(items))

	err = r.SaveItems(items)
	if err != nil {
		return fmt.Errorf("failed to save items: %w", err)
	}

	slog.Info("Marked todo as done", "index", i)
	return nil
}

// processDone is a wrapper for executeUpdate to satisfy the expected interface if needed.
func processDone(repo *todo.Repository, argStr string) error {
	return executeUpdate(repo, argStr)
}

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
