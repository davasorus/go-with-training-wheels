/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// renderList outputs the list of items to the provided writer.
func renderList(w io.Writer, items []todo.Todo, showOnlyDone bool) {
	tw := tabwriter.NewWriter(w, 3, 0, 1, ' ', 0)
	for _, i := range items {
		if !showOnlyDone || i.Done {
			fmt.Fprintln(tw, i.Label(), i.PrettyP()+"\t"+i.Text+"\t"+i.DoneStatus()+"\t")
		}
	}
	tw.Flush()
}

// listCmd returns the command to list tasks.
func listCmd(repo todo.TodoStore) *cobra.Command {
	var doneOpt bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks.",
		Long: `Display a table of tasks. 
Example: tri list --done`,
		Run: func(cmd *cobra.Command, args []string) {
			items, err := repo.ListItems()
			if err != nil {
				slog.Error("Error loading items", "err", err)
				return
			}

			sort.Sort(todo.ByPri(items))

			renderList(os.Stdout, items, doneOpt)
		},
	}

	cmd.Flags().BoolVar(&doneOpt, "done", false, "Show only completed items")
	return cmd
}
