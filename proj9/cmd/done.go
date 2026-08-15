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

// doneCmd represents the command to mark a task as complete.
var doneCmd = &cobra.Command{
	Use:     "done",
	Aliases: []string{"do"},
	Short:   "Mark a task as completed.",
	Long: `Mark a task as done using its position in the list.
Example: tri do 1`,
	Run: doneRun,
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

// doneRun handles the logic for marking a task as completed.
func doneRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		slog.Error("No index provided")
		return
	}

	err := processDone(args[0])
	if err != nil {
		slog.Error("Failed to mark done", "error", err)
		return
	}
}

// processDone handles the logic for finding and marking a task as completed.
func processDone(argStr string) error {
	items, err := todo.LoadItems()
	if err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}

	i, err := strconv.Atoi(argStr)
	if err != nil {
		return fmt.Errorf("invalid index provided: %w", err)
	}

	if i <= 0 || i > len(items) {
		return fmt.Errorf("%d does not match any item", i)
	}

	items[i-1].Done = true
	sort.Sort(todo.ByPri(items))

	err = todo.SaveItems(items)
	if err != nil {
		return fmt.Errorf("failed to save items: %w", err)
	}

	slog.Info("Marked todo as done", "index", i)
	return nil
}
