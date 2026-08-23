/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var text string
			err := WaitSpinner("Marking as done", func() error {
				var err error
				text, err = executeUpdate(repo, args[0], true)
				return err
			})
			if err != nil {
				return err
			}
			fmt.Printf("Marked done: %q\n", text)
			return nil
		},
	}
}

// undoCmd returns the command to mark a task as not completed.
func undoCmd(repo todo.TodoStore) *cobra.Command {
	return &cobra.Command{
		Use:     "undo",
		Aliases: []string{"undone"},
		Short:   "Mark a task as not completed.",
		Long: `Mark a done task as not done, using its number from 'tri list'.
Example: tri undo 1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var text string
			err := WaitSpinner("Marking as not done", func() error {
				var err error
				text, err = executeUpdate(repo, args[0], false)
				return err
			})
			if err != nil {
				return err
			}
			fmt.Printf("Marked not done: %q\n", text)
			return nil
		},
	}
}

// executeUpdate finds a task by its list position, sets its done state,
// and returns the task text for the success message.
func executeUpdate(r todo.TodoStore, argStr string, done bool) (string, error) {
	items, err := r.ListItems()
	if err != nil {
		return "", fmt.Errorf("failed to load items: %w", err)
	}

	// Match the ordering the user saw in `tri list`.
	sort.Sort(todo.ByPri(items))

	idx, err := validateIndex(argStr, len(items))
	if err != nil {
		return "", err
	}

	item := items[idx-1]
	if err := r.UpdateItemStatus(item.ID, done); err != nil {
		return "", fmt.Errorf("failed to update item: %w", err)
	}
	return item.Text, nil
}
