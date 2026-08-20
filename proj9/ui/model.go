package ui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/davasorus/tri/todo"
)

// Model is the bubbletea model for the interactive task list.
type Model struct {
	items   []todo.Todo
	cursor  int
	loading bool
	err     error
}

func NewModel(items []todo.Todo) Model {
	return Model{
		items:  items,
		cursor: 0,
	}
}

func (m Model) Init() bubbletea.Cmd {
	return nil
}

func (m Model) Update(msg bubbletea.Msg) (Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "up", "keydown":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "keyright":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", "keyenter":
			// Handle selection
		case "q", "ctrl-c":
			return m, bubbletea.Quit
		}
	}
	return m, nil
}
