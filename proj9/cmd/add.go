/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// priorityFlag is a pflag.Value that accepts numeric (0-2) and named
// (low, medium, high) priority values.
type priorityFlag int

func (p *priorityFlag) String() string {
	switch int(*p) {
	case 2:
		return "high"
	case 1:
		return "medium"
	default:
		return "low"
	}
}

func (p *priorityFlag) Set(s string) error {
	switch strings.ToLower(s) {
	case "0", "low":
		*p = 0
	case "1", "medium", "med":
		*p = 1
	case "2", "high":
		*p = 2
	default:
		return fmt.Errorf(`invalid priority %q: use low, medium, high, or 0-2`, s)
	}
	return nil
}

func (p *priorityFlag) Type() string { return "priority" }

// addCmd returns the command for adding a task.
func addCmd(repo todo.TodoStore) *cobra.Command {
	var priority priorityFlag
	var due string
	var tags []string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task with priority.",
		Long: `Add a new task to the system.
Example: tri add -p high --due tomorrow "Fix the build"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dueDate, err := parseDueDate(due)
			if err != nil {
				return err
			}

			items := prepareAddItems(args, int(priority), dueDate, tags)

			err = WaitSpinner("Saving items", func() error {
				return repo.SaveItems(items)
			})
			if err != nil {
				return fmt.Errorf("failed to save tasks: %w", err)
			}

			if len(items) == 1 {
				fmt.Printf("Added: %q\n", items[0].Text)
			} else {
				fmt.Printf("Added %d tasks.\n", len(items))
			}
			return nil
		},
	}
	cmd.Flags().VarP(&priority, "priority", "p", "Priority of the task: low, medium, or high (or 0-2).")
	cmd.Flags().StringVarP(&due, "due", "d", "", `Due date: YYYY-MM-DD, "today", or "tomorrow".`)
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tag for the task. Repeat the flag for more tags.")
	return cmd
}

// parseDueDate converts the --due flag value into a date.
// An empty value returns nil (no due date).
func parseDueDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch strings.ToLower(s) {
	case "today":
		return &today, nil
	case "tomorrow":
		t := today.AddDate(0, 0, 1)
		return &t, nil
	}

	t, err := time.ParseInLocation("2006-01-02", s, now.Location())
	if err != nil {
		return nil, fmt.Errorf(`invalid due date %q: use YYYY-MM-DD, "today", or "tomorrow"`, s)
	}
	return &t, nil
}

func prepareAddItems(args []string, priority int, dueDate *time.Time, tags []string) []todo.Todo {
	items := make([]todo.Todo, 0, len(args))
	for i, arg := range args {
		item := todo.Todo{Text: arg, DueDate: dueDate, Tags: tags}
		item.SetPriority(priority)
		item.Position = i + 1
		items = append(items, item)
	}
	return items
}
