/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// exportCmd returns the command to write all tasks as JSON to stdout.
func exportCmd(repo todo.TodoStore) *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Write all tasks as JSON to standard output.",
		Long: `Write all tasks as a JSON array to standard output, sorted the
same way as 'tri list'. Useful for scripting:
Example: tri export | jq '.[] | select(.done)'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := repo.ListItems()
			if err != nil {
				return fmt.Errorf("failed to load items: %w", err)
			}

			sort.Sort(todo.ByPri(items))

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(items)
		},
	}
}
