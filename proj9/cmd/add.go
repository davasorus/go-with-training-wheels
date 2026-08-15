/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// addCmd represents the command to add a task.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task with priority.",
	Long: `Add a new task to the system. 
Example: tri add -p high MyTask`,
	Run: addRun,
}

var priority int

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().IntVarP(&priority, "priority", "p", 0, "Priority of the task (0=Low, 1=Medium, 2=High).")
}

// addRun handles the logic for adding a new task.
func addRun(cmd *cobra.Command, args []string) {

	items := []todo.Todo{}
	for i, arg := range args {
		item := todo.Todo{Text: arg}
		item.SetPriority(priority)
		item.Position = i + 1
		items = append(items, item)
	}

	err := todo.SaveItems(items)
	if err != nil {
		fmt.Println("Error saving items:", err)
	}

}
