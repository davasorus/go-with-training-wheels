/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// listCmd represents the command to list tasks.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks.",
	Long: `Display a table of tasks. 
Example: tri list --done`,
	Run: writeOutput,
}

var (
	doneOpt bool
	allOpt  bool
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&doneOpt, "done", false, "Show only completed items")
}

// writeOutput displays the list of tasks in a table format.
func writeOutput(cmd *cobra.Command, args []string) {
	items, err := todo.LoadItems()
	if err != nil {
		slog.Error("Error loading items", "err", err)
		return
	}

	sort.Sort(todo.ByPri(items))

	w := tabwriter.NewWriter(os.Stdout, 3, 0, 1, ' ', 0)

	for _, i := range items {
		if !doneOpt || i.Done {
			fmt.Fprintln(w, i.Label(), i.PrettyP()+"\t"+i.Text+"\t"+i.DoneStatus()+"\t")
		}
	}

	w.Flush()
}
