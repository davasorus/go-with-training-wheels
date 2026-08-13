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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: writeOutput,
}

var (
	doneOpt bool
	allOpt  bool
)

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	listCmd.Flags().BoolVar(&doneOpt, "done", false, "Show only done items")
}

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
