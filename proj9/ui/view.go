package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	helpStyle    = lipgloss.NewStyle().Faint(true)
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	overdueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	tagStyle     = lipgloss.NewStyle().Faint(true)
	selectionStr = " > "
)

func (m Model) View() string {
	var s string
	s += titleStyle.Render("Tasks") + "\n\n"

	if len(m.items) == 0 {
		s += "No tasks.\n"
	}

	for i, item := range m.items {
		label := "   "
		if i == m.cursor {
			label = selectionStr
		}

		status := "[ ]"
		if item.Done {
			status = "[x]"
		}

		// High priority renders red; other items render default.
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
		if item.Priority == 2 {
			style = style.Foreground(lipgloss.Color("1"))
		}

		line := fmt.Sprintf("%s%d. %s %s", label, i+1, status, style.Render(item.Text))
		if len(item.Tags) > 0 {
			line += " " + tagStyle.Render("#"+strings.Join(item.Tags, " #"))
		}
		if item.DueDate != nil {
			due := item.DueDate.Format("2006-01-02")
			if !item.Done && item.DueDate.Before(time.Now()) {
				line += " " + overdueStyle.Render(due+" (overdue)")
			} else {
				line += " " + helpStyle.Render(due)
			}
		}
		s += line + "\n"
	}

	if m.adding {
		s += "\nAdd task: " + m.input.View() + "\n"
	}

	if m.loading {
		s += "\nWorking...\n"
	}

	if m.err != nil {
		s += "\n" + errStyle.Render("Error: "+m.err.Error()) + "\n"
	}

	if m.adding {
		s += helpStyle.Render("\nenter save · esc cancel") + "\n"
	} else {
		s += helpStyle.Render("\n↑/k up · ↓/j down · enter/space toggle done · a add · d delete · q quit") + "\n"
	}

	return s
}
