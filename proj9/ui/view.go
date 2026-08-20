package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/davasorus/tri/todo"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	selectionStr = " > "
)

func Render(m Model) string {
	var s string
	s += titleStyle.Render("Select a task to focus on:") + "\n"

	for i, item := range m.items {
		label := ""
		if i == m.cursor {
			label = selectionStr
		}

		// Apply consistent styling for high priority (red) and done items (strike-through)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
		if item.Priority == 2 {
			style = style.Foreground(lipgloss.Color("1"))
		}

		text := item.Text
		if item.Done {
			text = fmt.Sprintf("~~%s~~", text)
		}

		s += fmt.Sprintf("%s %d. %s\n", label, i+1, style.Render(text))
	}

	if m.err != nil {
		s += "\nError: " + m.err.Error()
	}

	return s
}
