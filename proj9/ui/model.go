package ui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/davasorus/tri/todo"
)

// mode is the input mode of the interactive view.
type mode int

const (
	modeNormal mode = iota
	modeAdd
	modeEdit
)

// refreshMsg carries the result of a repository operation back into Update.
type refreshMsg struct {
	items []todo.Todo
	err   error
	// followID keeps the cursor on this task after the list changes.
	// Zero means: keep the cursor at its clamped index.
	followID int
}

// Model is the bubbletea model for the interactive task list.
type Model struct {
	repo    todo.TodoStore
	items   []todo.Todo
	filter  func(todo.Todo) bool
	cursor  int
	offset  int // first visible row for scrolling
	height  int // terminal height from the last WindowSizeMsg
	loading bool
	mode    mode
	editID  int // task being edited in modeEdit
	input   textinput.Model
	err     error
}

// NewModel builds the model. filter may be nil, in which case all
// items are shown. A non-nil filter is re-applied on every refresh so
// command-line filters survive toggles and deletes.
func NewModel(repo todo.TodoStore, items []todo.Todo, filter func(todo.Todo) bool) Model {
	input := textinput.New()
	input.Placeholder = "New task text"
	input.CharLimit = 200

	return Model{
		repo:   repo,
		items:  items,
		filter: filter,
		input:  input,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// visibleRows is how many task lines fit in the current terminal.
func (m Model) visibleRows() int {
	// Title, blank line, optional input line, status lines, help line.
	const chrome = 7
	if m.height <= 0 {
		return 20
	}
	rows := m.height - chrome
	if rows < 3 {
		return 3
	}
	return rows
}

// clampView keeps the cursor inside the list and the window around
// the cursor.
func (m *Model) clampView() {
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	rows := m.visibleRows()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+rows {
		m.offset = m.cursor - rows + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

// refresh reloads all items from the repository, sorted the same way
// the non-interactive list is.
func (m Model) refresh(followID int) tea.Msg {
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
	return refreshMsg{items: items, followID: followID}
}

// toggleCurrent flips the done state of the item under the cursor.
func (m Model) toggleCurrent() tea.Cmd {
	item := m.items[m.cursor]
	return func() tea.Msg {
		if err := m.repo.UpdateItemStatus(item.ID, !item.Done); err != nil {
			return refreshMsg{items: m.items, err: err}
		}
		return m.refresh(item.ID)
	}
}

// setPriority changes the priority of the item under the cursor.
func (m Model) setPriority(priority int) tea.Cmd {
	item := m.items[m.cursor]
	item.SetPriority(priority)
	return func() tea.Msg {
		if err := m.repo.UpdateItem(item); err != nil {
			return refreshMsg{items: m.items, err: err}
		}
		return m.refresh(item.ID)
	}
}

// addItem saves a new task with the given text.
func (m Model) addItem(text string) tea.Cmd {
	return func() tea.Msg {
		item := todo.Todo{Text: text, Position: len(m.items) + 1}
		if err := m.repo.SaveItems([]todo.Todo{item}); err != nil {
			return refreshMsg{items: m.items, err: err}
		}
		return m.refresh(0)
	}
}

// editItem saves new text for the task that was being edited.
func (m Model) editItem(id int, text string) tea.Cmd {
	var item todo.Todo
	for _, it := range m.items {
		if it.ID == id {
			item = it
			break
		}
	}
	item.Text = text
	return func() tea.Msg {
		if err := m.repo.UpdateItem(item); err != nil {
			return refreshMsg{items: m.items, err: err}
		}
		return m.refresh(id)
	}
}

// deleteCurrent removes the item under the cursor.
func (m Model) deleteCurrent() tea.Cmd {
	item := m.items[m.cursor]
	return func() tea.Msg {
		if err := m.repo.DeleteItem(item.ID); err != nil {
			return refreshMsg{items: m.items, err: err}
		}
		return m.refresh(0)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.clampView()
		return m, nil

	case refreshMsg:
		m.loading = false
		m.err = msg.err
		m.items = msg.items
		if msg.followID != 0 {
			for i, it := range m.items {
				if it.ID == msg.followID {
					m.cursor = i
					break
				}
			}
		}
		m.clampView()
		return m, nil

	case tea.KeyMsg:
		// Input modes capture all keys except their own controls.
		if m.mode == modeAdd || m.mode == modeEdit {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.mode = modeNormal
				m.input.Reset()
				return m, nil
			case "enter":
				text := strings.TrimSpace(m.input.Value())
				current := m.mode
				m.mode = modeNormal
				m.input.Reset()
				if text == "" {
					return m, nil
				}
				m.loading = true
				m.err = nil
				if current == modeEdit {
					return m, m.editItem(m.editID, text)
				}
				return m, m.addItem(text)
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

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
			m.clampView()
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
			m.clampView()
		case "enter", " ":
			if len(m.items) == 0 {
				return m, nil
			}
			m.loading = true
			m.err = nil
			return m, m.toggleCurrent()
		case "1", "2", "3":
			if len(m.items) == 0 {
				return m, nil
			}
			m.loading = true
			m.err = nil
			return m, m.setPriority(int(msg.String()[0] - '1'))
		case "d":
			if len(m.items) == 0 {
				return m, nil
			}
			m.loading = true
			m.err = nil
			return m, m.deleteCurrent()
		case "a":
			m.mode = modeAdd
			m.err = nil
			m.input.Placeholder = "New task text"
			m.input.Focus()
			return m, textinput.Blink
		case "e":
			if len(m.items) == 0 {
				return m, nil
			}
			item := m.items[m.cursor]
			m.mode = modeEdit
			m.editID = item.ID
			m.err = nil
			m.input.Placeholder = "Task text"
			m.input.SetValue(item.Text)
			m.input.CursorEnd()
			m.input.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}
