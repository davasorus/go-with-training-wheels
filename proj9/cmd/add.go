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

// addCmd returns the command for adding a task.
func addCmd(repo todo.TodoStore) *cobra.Command {
	var priority int
	var due string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task with priority.",
		Long: `Add a new task to the system.
Example: tri add -p 2 --due tomorrow "Fix the build"`,
		Run: func(cmd *cobra.Command, args []string) {
			dueDate, err := parseDueDate(due)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			items := prepareAddItems(args, priority, dueDate)

			err = WaitSpinner("Saving items", func() error {
				return repo.SaveItems(items)
			})

			if err != nil {
				fmt.Println("Error: Failed to save tasks. Please check your input and try again.")
			}
		},
	}
	cmd.Flags().IntVarP(&priority, "priority", "p", 0, "Priority of the task (0=Low, 1=Medium, 2=High).")
	cmd.Flags().StringVarP(&due, "due", "d", "", `Due date: YYYY-MM-DD, "today", or "tomorrow".`)
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

func prepareAddItems(args []string, priority int, dueDate *time.Time) []todo.Todo {
	items := make([]todo.Todo, 0, len(args))
	for i, arg := range args {
		item := todo.Todo{Text: arg, DueDate: dueDate}
		item.SetPriority(priority)
		item.Position = i + 1
		items = append(items, item)
	}
	return items
}
