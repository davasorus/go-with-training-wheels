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
	_, _ = fmt.Fprintln(tw, "ID", "Priority", "Task", "Tags", "Status", "Due Date")
	_, _ = fmt.Fprintln(tw)
	now := time.Now()
	for idx, i := range items {
		if !showOnlyDone || i.Done {
			dateStr := ""
			if i.DueDate != nil {
				dateStr = i.DueDate.Format("2006-01-02")
				if !i.Done && i.DueDate.Before(now) {
					dateStr = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(dateStr + " (overdue)")
				}
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

			taskLabel := fmt.Sprintf("%d.", idx+1)
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

			tagStr := strings.Join(i.Tags, ",")
			_, _ = fmt.Fprintln(tw, style.Render(taskLabel), i.PrettyP(), i.Text, tagStr, statusStyle.Render(status), dateStr)
		}
	}
	if err := tw.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering list: %v\n", err)
	}
}

// buildFilter composes the list filters into one predicate.
// It returns nil when no filter flags are active.
func buildFilter(doneOpt bool, nearOpt int, queryOpt string, tagOpts []string) func(todo.Todo) bool {
	if !doneOpt && nearOpt <= 0 && queryOpt == "" && len(tagOpts) == 0 {
		return nil
	}
	now := time.Now()
	return func(itm todo.Todo) bool {
		if doneOpt && !itm.Done {
			return false
		}
		if nearOpt > 0 {
			if itm.DueDate == nil {
				return false
			}
			daysUntil := itm.DueDate.Sub(now).Hours() / 24
			if daysUntil < 0 || daysUntil > float64(nearOpt) {
				return false
			}
		}
		if queryOpt != "" {
			q := strings.ToLower(queryOpt)
			if !strings.Contains(strings.ToLower(itm.Label()), q) &&
				!strings.Contains(strings.ToLower(itm.Text), q) {
				return false
			}
		}
		for _, tag := range tagOpts {
			if !itm.HasTag(tag) {
				return false
			}
		}
		return true
	}
}

func applyFilter(items []todo.Todo, filter func(todo.Todo) bool) []todo.Todo {
	if filter == nil {
		return items
	}
	var kept []todo.Todo
	for _, itm := range items {
		if filter(itm) {
			kept = append(kept, itm)
		}
	}
	return kept
}

func runTUI(repo todo.TodoStore, items []todo.Todo, filter func(todo.Todo) bool) error {
	m := ui.NewModel(repo, items, filter)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run the interactive view: %w", err)
	}
	return nil
}

// listCmd returns the command to list tasks.
func listCmd(repo todo.TodoStore) *cobra.Command {
	var queryOpt string
	var doneOpt bool
	var nearOpt int
	var interactiveOpt bool
	var tagOpts []string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks.",
		Long: `Display a table of tasks. 
Example: tri list --done`,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := repo.ListItems()
			if err != nil {
				if strings.Contains(err.Error(), "does not exist") {
					return fmt.Errorf("%w\nHint: the database schema may be out of date. Run 'tri migrateDb'", err)
				}
				return err
			}

			sort.Sort(todo.ByPri(items))

			filter := buildFilter(doneOpt, nearOpt, queryOpt, tagOpts)
			filteredItems := applyFilter(items, filter)

			if interactiveOpt {
				return runTUI(repo, filteredItems, filter)
			}

			renderList(os.Stdout, filteredItems, doneOpt, nearOpt, queryOpt)
			return nil
		},
	}

	cmd.Flags().BoolVar(&doneOpt, "done", false, "Show only completed items")
	cmd.Flags().IntVarP(&nearOpt, "near", "n", 0, "Show tasks due within the next 7 days (default: all)")
	cmd.Flags().StringVarP(&queryOpt, "query", "q", "", "Filter items by title or description")
	cmd.Flags().BoolVarP(&interactiveOpt, "interactive", "i", false, "Start in interactive mode (TUI)")
	cmd.Flags().StringSliceVar(&tagOpts, "tag", nil, "Show only tasks with this tag. Repeat the flag to require more tags.")
	return cmd
}
