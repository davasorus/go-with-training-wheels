/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"log/slog"
	"sort"
	"strconv"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// doneCmd represents the command to mark a task as complete.
var doneCmd = &cobra.Command{
	Use:     "done",
	Aliases: []string{"do"},
	Short:   "Mark a task as completed.",
	Long:    `Mark a task as done using its position in the list.
Example: tri do 1`,
	Run: doneRun,
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

// doneRun handles the logic for marking a task as completed.
func doneRun(cmd *cobra.Command, args []string) {
	items, err := todo.LoadItems()
	if err != nil {
		slog.Error("Failed to load items", "err", err)
		return
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		slog.Error("Invalid index provided", "err", err)
		return
	}

	if i > 0 && i <= len(items) {
		items[i-1].Done = true
		slog.Info("Marked todo as done", "todo", items[i-1])

		sort.Sort(todo.ByPri(items))

		err = todo.SaveItems(items)
		if err != nil {
			slog.Error("Failed to save items", "err", err)
			return
		}
		slog.Debug("Items saved successfully", "count", len(items))
	} else {
		log.Println(i, "does not match any item")
	}
}
