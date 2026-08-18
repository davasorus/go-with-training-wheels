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

			tw := tabwriter.NewWriter(os.Stdout, 3, 0, 1, ' ', 0)
			for _, i := range items {
				if !doneOpt || i.Done {
					fmt.Fprintln(tw, i.Label(), i.PrettyP()+"\t"+i.Text+"\t"+i.DoneStatus()+"\t")
				}
			}
			tw.Flush()
		},
	}

	cmd.Flags().BoolVar(&doneOpt, "done", false, "Show only completed items")
	return cmd
}
