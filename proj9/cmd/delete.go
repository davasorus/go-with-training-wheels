/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
	var yes bool
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "Remove a task from the list.",
		Long: `Remove a task from the list based on its position or ID.
Example: tri delete 1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			text, id, err := resolveDeleteTarget(repo, args[0])
			if err != nil {
				return err
			}

			if !yes && !confirm(fmt.Sprintf("Delete %q?", text)) {
				fmt.Println("Canceled.")
				return nil
			}

			err = WaitSpinner("Deleting item", func() error {
				return repo.DeleteItem(id)
			})
			if err != nil {
				return err
			}
			fmt.Printf("Deleted: %q\n", text)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Delete without asking for confirmation.")
	return cmd
}

// resolveDeleteTarget maps the argument to a task. A number within the
// list range is a list position; a larger number is a raw database ID.
func resolveDeleteTarget(r todo.TodoStore, argStr string) (string, int, error) {
	val, err := strconv.Atoi(argStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid input: %s is not a number", argStr)
	}

	items, err := r.ListItems()
	if err != nil {
		return "", 0, fmt.Errorf("failed to load items: %w", err)
	}

	// Match the ordering the user saw in `tri list`.
	sort.Sort(todo.ByPri(items))

	if val >= 1 && val <= len(items) {
		item := items[val-1]
		return item.Text, item.ID, nil
	}

	// Fall back to treating the value as a raw database ID.
	for _, item := range items {
		if item.ID == val {
			return item.Text, item.ID, nil
		}
	}
	return "", 0, fmt.Errorf("%d does not match any item", val)
}
