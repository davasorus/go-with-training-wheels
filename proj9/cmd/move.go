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

// moveCmd returns the command to reorder a task within the list.
func moveCmd(repo todo.TodoStore) *cobra.Command {
	return &cobra.Command{
		Use:   "move",
		Short: "Move a task up or down in the list.",
		Long: `Move a task one place up or down. The first argument is the
number from 'tri list'. The second argument is "up" or "down".
Movement is within the same priority and done state.
Example: tri move 3 up`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := WaitSpinner("Moving item", func() error {
				return executeMove(repo, args[0], args[1])
			})
			if err != nil {
				return err
			}
			fmt.Printf("Moved task %s %s.\n", args[0], args[1])
			return nil
		},
	}
}

// executeMove swaps a task with its neighbor in the displayed order.
func executeMove(r todo.TodoStore, argStr, direction string) error {
	if direction != "up" && direction != "down" {
		return fmt.Errorf(`invalid direction %q: use "up" or "down"`, direction)
	}

	items, err := r.ListItems()
	if err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}

	sort.Sort(todo.ByPri(items))

	idx, err := validateIndex(argStr, len(items))
	if err != nil {
		return err
	}

	target := idx - 1
	neighbor := target - 1
	if direction == "down" {
		neighbor = target + 1
	}
	if neighbor < 0 || neighbor >= len(items) {
		return fmt.Errorf("task %d is already at the %s of the list", idx, map[string]string{"up": "top", "down": "bottom"}[direction])
	}

	a, b := items[target], items[neighbor]
	if a.Done != b.Done || a.Priority != b.Priority {
		return fmt.Errorf("cannot move across a priority or done boundary: change the priority with 'tri edit %d -p ...' instead", idx)
	}

	// Positions may be stale duplicates from older versions. Normalize
	// every position to the current displayed order, then swap.
	for i := range items {
		items[i].Position = i + 1
	}
	items[target].Position, items[neighbor].Position = items[neighbor].Position, items[target].Position

	for _, it := range []todo.Todo{items[target], items[neighbor]} {
		if err := r.UpdateItem(it); err != nil {
			return err
		}
	}

	// Persist the normalized positions for the rest so future moves
	// stay stable.
	for i, it := range items {
		if i == target || i == neighbor {
			continue
		}
		if err := r.UpdateItem(it); err != nil {
			return err
		}
	}
	return nil
}
