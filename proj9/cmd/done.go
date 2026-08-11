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

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:     "done",
	Aliases: []string{"do"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: doneRun,
}

func init() {
	rootCmd.AddCommand(doneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func doneRun(cmd *cobra.Command, args []string) {
	items, err := todo.LoadItems(dataFile)
	i, err := strconv.Atoi(args[0])
	if err != nil {
		slog.Error("Error converting argument to integer", "err", err)
		return
	}

	if i > 0 && i <= len(items) {
		items[i-1].Done = true
		slog.Info("Marked item as done", "item", items[i-1])

		sort.Sort(todo.ByPri(items))

		err = todo.SaveItems(dataFile, items)
		if err != nil {
			slog.Error("Error saving items", "err", err)
			return
		}
		slog.Debug("Items saved successfully", "count", len(items))
	} else {
		log.Println(i, "doesn't match any items")
	}

}
