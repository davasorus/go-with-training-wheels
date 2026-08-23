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

// editCmd returns the command to change an existing task.
func editCmd(repo todo.TodoStore) *cobra.Command {
	var priority priorityFlag = -1
	var text, due string
	var tags []string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Change the text, priority, due date, or tags of a task.",
		Long: `Change an existing task. The argument is the number from 'tri list'.
Only the flags you provide are changed.
Example: tri edit 2 --text "New description" -p high --due 2026-09-01
Use --due none to remove the due date. Use --tags none to remove all tags.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Error: No number provided. Example: tri edit 1 --text \"New text\"")
				return
			}

			err := WaitSpinner("Updating item", func() error {
				return executeEdit(repo, args[0], cmd.Flags().Changed("priority"), int(priority), text, due, tags)
			})

			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		},
	}
	cmd.Flags().VarP(&priority, "priority", "p", "New priority: low, medium, or high (or 0-2).")
	cmd.Flags().StringVarP(&text, "text", "t", "", "New task text.")
	cmd.Flags().StringVarP(&due, "due", "d", "", `New due date: YYYY-MM-DD, "today", "tomorrow", or "none" to clear.`)
	cmd.Flags().StringSliceVar(&tags, "tags", nil, `New tags (replaces existing). Use "none" to clear.`)
	return cmd
}

// executeEdit applies the requested changes to the task at the given
// list position.
func executeEdit(r todo.TodoStore, argStr string, priorityChanged bool, priority int, text, due string, tags []string) error {
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
	item := items[idx-1]

	changed := false

	if text != "" {
		item.Text = text
		changed = true
	}

	if priorityChanged {
		item.SetPriority(priority)
		changed = true
	}

	if due != "" {
		if due == "none" {
			item.DueDate = nil
		} else {
			d, err := parseDueDate(due)
			if err != nil {
				return err
			}
			item.DueDate = d
		}
		changed = true
	}

	if tags != nil {
		if len(tags) == 1 && tags[0] == "none" {
			item.Tags = nil
		} else {
			item.Tags = tags
		}
		changed = true
	}

	if !changed {
		return fmt.Errorf("nothing to change: provide at least one of --text, --priority, --due, --tags")
	}

	return r.UpdateItem(item)
}
