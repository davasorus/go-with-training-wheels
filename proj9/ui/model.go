package ui

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/davasorus/tri/todo"
)

// refreshMsg carries the result of a repository operation back into Update.
type refreshMsg struct {
	items []todo.Todo
	err   error
}

// Model is the bubbletea model for the interactive task list.
type Model struct {
	repo    todo.TodoStore
	items   []todo.Todo
	filter  func(todo.Todo) bool
	cursor  int
	loading bool
	err     error
}

// NewModel builds the model. filter may be nil, in which case all
// items are shown. A non-nil filter is re-applied on every refresh so
// command-line filters survive toggles and deletes.
func NewModel(repo todo.TodoStore, items []todo.Todo, filter func(todo.Todo) bool) Model {
	return Model{
		repo:   repo,
		items:  items,
		filter: filter,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// refresh reloads all items from the repository, sorted the same way
// the non-interactive list is.
func (m Model) refresh() tea.Msg {
	items, err := m.repo.ListItems()
	if err != nil {
		return refreshMsg{items: m.items, err: err}
	}
	if m.filter != nil {
		kept := items[:0]
		for _, it := range items {
			if m.filter(it) {
				kept = append(kept, it)
			}
		}
		items = kept
	}
	sort.Sort(todo.ByPri(items))
	return refreshMsg{items: items}
}

// toggleCurrent flips the done state of the item under the cursor.
func (m Model) toggleCurrent() tea.Cmd {
	item := m.items[m.cursor]
	return func() tea.Msg {
		if err := m.repo.UpdateItemStatus(item.ID, !item.Done); err != nil {
			return refreshMsg{items: m.items, err: err}
		}
		return m.refresh()
	}
}

// deleteCurrent removes the item under the cursor.
func (m Model) deleteCurrent() tea.Cmd {
	item := m.items[m.cursor]
	return func() tea.Msg {
		if err := m.repo.DeleteItem(item.ID); err != nil {
			return refreshMsg{items: m.items, err: err}
		}
		return m.refresh()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		m.loading = false
		m.err = msg.err
		m.items = msg.items
		if m.cursor >= len(m.items) {
			m.cursor = len(m.items) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case tea.KeyMsg:
		// Always allow quitting, even mid-operation.
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		// Ignore other input while an operation is in flight.
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.items) == 0 {
				return m, nil
			}
			m.loading = true
			m.err = nil
			return m, m.toggleCurrent()
		case "d":
			if len(m.items) == 0 {
				return m, nil
			}
			m.loading = true
			m.err = nil
			return m, m.deleteCurrent()
		}
	}
	return m, nil
}
