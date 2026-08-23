/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davasorus/tri/todo"
	"github.com/davasorus/tri/ui"
	"github.com/spf13/cobra"
)

// renderList outputs the list of items to the provided writer.
func renderList(w io.Writer, items []todo.Todo, showOnlyDone bool, nearDays int, queryOpt string) {
	tw := tabwriter.NewWriter(w, 3, 0, 1, ' ', 0)
	_, _ = fmt.Fprintln(tw, "ID", "Priority", "Task", "Status", "Due Date")
	_, _ = fmt.Fprintln(tw)
	now := time.Now()
	for _, i := range items {
		if !showOnlyDone || i.Done {
			dateStr := ""
			if i.DueDate != nil {
				dateStr = i.DueDate.Format("2006-01-02")
				if nearDays > 0 {
					daysUntil := i.DueDate.Sub(now).Hours() / 24
					if daysUntil < 0 || daysUntil > float64(nearDays) {
						continue
					}
				}
			}

			if queryOpt != "" {
				found := false
				if strings.Contains(strings.ToLower(i.Label()), strings.ToLower(queryOpt)) ||
					strings.Contains(strings.ToLower(i.Text), strings.ToLower(queryOpt)) {
					found = true
				}
				if !found {
					continue
				}
			}

			var style = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
			if i.Priority == 2 {
				style = style.Foreground(lipgloss.Color("1"))
			}

			taskLabel := i.Label()
			if i.Done {
				taskLabel = style.Render(fmt.Sprintf("~~%s~~", taskLabel))
			} else {
				taskLabel = style.Render(taskLabel)
			}

			status := "[ ]"
			statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
			if i.Done {
				status = "[x]"
				statusStyle = statusStyle.Foreground(lipgloss.Color("3"))
			}

			_, _ = fmt.Fprintln(tw, style.Render(taskLabel), i.PrettyP(), i.Text, statusStyle.Render(status), dateStr)
		}
	}
	if err := tw.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering list: %v\n", err)
	}
}

func runTUI(repo todo.TodoStore, items []todo.Todo) {
	m := ui.NewModel(repo, items)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
	}
}

// listCmd returns the command to list tasks.
func listCmd(repo todo.TodoStore) *cobra.Command {
	var queryOpt string
	var doneOpt bool
	var nearOpt int
	var interactiveOpt bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks.",
		Long: `Display a table of tasks. 
Example: tri list --done`,
		Run: func(cmd *cobra.Command, args []string) {
			items, err := repo.ListItems()
			if err != nil {
				fmt.Println("Error: Failed to load items from the database.")
				return
			}

			sort.Sort(todo.ByPri(items))

			// Apply filters to the items list
			filteredItems := items
			if doneOpt {
				var filtered []todo.Todo
				for _, itm := range items {
					if itm.Done {
						filtered = append(filtered, itm)
					}
				}
				filteredItems = filtered
			}

			if nearOpt > 0 {
				var filtered []todo.Todo
				now := time.Now()
				for _, itm := range items {
					if itm.DueDate != nil {
						daysUntil := itm.DueDate.Sub(now).Hours() / 24
						if daysUntil >= 0 && daysUntil <= float64(nearOpt) {
							filtered = append(filtered, itm)
						}
					}
				}
				filteredItems = filtered
			}

			if queryOpt != "" {
				var filtered []todo.Todo
				for _, itm := range items {
					if strings.Contains(strings.ToLower(itm.Label()), strings.ToLower(queryOpt)) ||
						strings.Contains(strings.ToLower(itm.Text), strings.ToLower(queryOpt)) {
						filtered = append(filtered, itm)
					}
				}
				filteredItems = filtered
			}

			if interactiveOpt {
				runTUI(repo, filteredItems)
				return
			}

			renderList(os.Stdout, filteredItems, doneOpt, nearOpt, queryOpt)
		},
	}

	cmd.Flags().BoolVar(&doneOpt, "done", false, "Show only completed items")
	cmd.Flags().IntVarP(&nearOpt, "near", "n", 0, "Show tasks due within the next 7 days (default: all)")
	cmd.Flags().StringVarP(&queryOpt, "query", "q", "", "Filter items by title or description")
	cmd.Flags().BoolVar(&interactiveOpt, "interactive", false, "Start in interactive mode (TUI)")
	return cmd
}
