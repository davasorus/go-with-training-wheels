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

	rows := m.visibleRows()
	end := m.offset + rows
	if end > len(m.items) {
		end = len(m.items)
	}

	if m.offset > 0 {
		s += helpStyle.Render(fmt.Sprintf("  ↑ %d more", m.offset)) + "\n"
	}

	for i := m.offset; i < end; i++ {
		item := m.items[i]
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

	if end < len(m.items) {
		s += helpStyle.Render(fmt.Sprintf("  ↓ %d more", len(m.items)-end)) + "\n"
	}

	switch m.mode {
	case modeAdd:
		s += "\nAdd task: " + m.input.View() + "\n"
	case modeEdit:
		s += "\nEdit task: " + m.input.View() + "\n"
	}

	if m.loading {
		s += "\nWorking...\n"
	}

	if m.err != nil {
		s += "\n" + errStyle.Render("Error: "+m.err.Error()) + "\n"
	}

	if m.mode == modeAdd || m.mode == modeEdit {
		s += helpStyle.Render("\nenter save · esc cancel") + "\n"
	} else {
		s += helpStyle.Render("\n↑/k up · ↓/j down · enter/space toggle · a add · e edit · 1/2/3 priority · d delete · q quit") + "\n"
	}

	return s
}
