/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// addCmd returns the command for adding a task.
func addCmd(repo todo.TodoStore) *cobra.Command {
	var priority int
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task with priority.",
		Long: `Add a new task to the system. 
Example: tri add -p high MyTask`,
		Run: func(cmd *cobra.Command, args []string) {
			items := prepareAddItems(args, priority)

			err := repo.SaveItems(items)
			if err != nil {
				fmt.Println("Error saving items:", err)
			}
		},
	}
	cmd.Flags().IntVarP(&priority, "priority", "p", 0, "Priority of the task (0=Low, 1=Medium, 2=High).")
	return cmd
}

func prepareAddItems(args []string, priority int) []todo.Todo {
	items := make([]todo.Todo, 0, len(args))
	for i, arg := range args {
		item := todo.Todo{Text: arg}
		item.SetPriority(priority)
		item.Position = i + 1
		items = append(items, item)
	}
	return items
}
