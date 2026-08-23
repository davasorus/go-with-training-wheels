package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	helpStyle    = lipgloss.NewStyle().Faint(true)
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
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

		s += fmt.Sprintf("%s%d. %s %s\n", label, i+1, status, style.Render(item.Text))
	}

	if m.loading {
		s += "\nWorking...\n"
	}

	if m.err != nil {
		s += "\n" + errStyle.Render("Error: "+m.err.Error()) + "\n"
	}

	s += helpStyle.Render("\n↑/k up · ↓/j down · enter/space toggle done · d delete · q quit") + "\n"

	return s
}
