package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/davasorus/tri/todo"
)

// Model is the bubbletea model for the interactive task list.
type Model struct {
	items   []todo.Todo
	cursor int
	err    error
}

func NewModel(items []todo.Todo) Model {
	return Model{
		items:  items,
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			// Handle selection
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
