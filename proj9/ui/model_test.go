package ui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davasorus/tri/todo"
)

// fakeStore is an in-memory TodoStore for tests.
type fakeStore struct {
	items     []todo.Todo
	listErr   error
	updateErr error
	deleteErr error
}

func (f *fakeStore) ListItems() ([]todo.Todo, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	out := make([]todo.Todo, len(f.items))
	copy(out, f.items)
	return out, nil
}

func (f *fakeStore) SaveItems(items []todo.Todo) error { return nil }

func (f *fakeStore) UpdateItemStatus(id int, done bool) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	for i := range f.items {
		if f.items[i].ID == id {
			f.items[i].Done = done
			return nil
		}
	}
	return errors.New("not found")
}

func (f *fakeStore) DeleteItem(id int) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	for i := range f.items {
		if f.items[i].ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func testItems() []todo.Todo {
	return []todo.Todo{
		{ID: 10, Text: "Task A", Priority: 2, Position: 1},
		{ID: 20, Text: "Task B", Priority: 1, Position: 2},
		{ID: 30, Text: "Task C", Priority: 0, Position: 3},
	}
}

func newTestModel(store *fakeStore, filter func(todo.Todo) bool) Model {
	items, _ := store.ListItems()
	return NewModel(store, items, filter)
}

func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "space":
		return tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

// update runs Update and casts the result back to Model.
func update(t *testing.T, m Model, msg tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(msg)
	nm, ok := next.(Model)
	require.True(t, ok, "Update must return a ui.Model")
	return nm, cmd
}

func TestCursorMovement(t *testing.T) {
	m := newTestModel(&fakeStore{items: testItems()}, nil)

	m, _ = update(t, m, key("down"))
	assert.Equal(t, 1, m.cursor)

	m, _ = update(t, m, key("j"))
	assert.Equal(t, 2, m.cursor)

	// Clamp at the bottom.
	m, _ = update(t, m, key("down"))
	assert.Equal(t, 2, m.cursor)

	m, _ = update(t, m, key("k"))
	assert.Equal(t, 1, m.cursor)

	m, _ = update(t, m, key("up"))
	assert.Equal(t, 0, m.cursor)

	// Clamp at the top.
	m, _ = update(t, m, key("up"))
	assert.Equal(t, 0, m.cursor)
}

func TestQuitKeys(t *testing.T) {
	for _, k := range []string{"q", "ctrl+c"} {
		m := newTestModel(&fakeStore{items: testItems()}, nil)
		_, cmd := update(t, m, key(k))
		require.NotNil(t, cmd, "key %q must return a command", k)
		assert.IsType(t, tea.QuitMsg{}, cmd(), "key %q must quit", k)
	}
}

func TestToggleDone(t *testing.T) {
	store := &fakeStore{items: testItems()}
	m := newTestModel(store, nil)

	m, cmd := update(t, m, key("enter"))
	assert.True(t, m.loading, "toggle must set loading")
	require.NotNil(t, cmd)

	msg := cmd()
	refresh, ok := msg.(refreshMsg)
	require.True(t, ok, "command must produce a refreshMsg")
	require.NoError(t, refresh.err)

	m, _ = update(t, m, refresh)
	assert.False(t, m.loading, "refresh must clear loading")

	// The cursor was on ID 10 (highest priority sorts first).
	var found bool
	for _, it := range m.items {
		if it.ID == 10 {
			found = true
			assert.True(t, it.Done, "item under cursor must be toggled")
		}
	}
	assert.True(t, found)
}

func TestSpaceAlsoToggles(t *testing.T) {
	m := newTestModel(&fakeStore{items: testItems()}, nil)
	m, cmd := update(t, m, key("space"))
	assert.True(t, m.loading)
	assert.NotNil(t, cmd)
}

func TestDeleteClampsCursor(t *testing.T) {
	store := &fakeStore{items: testItems()}
	m := newTestModel(store, nil)

	// Move to the last item and delete it.
	m.cursor = 2
	m, cmd := update(t, m, key("d"))
	require.NotNil(t, cmd)

	msg := cmd()
	refresh, ok := msg.(refreshMsg)
	require.True(t, ok)
	require.NoError(t, refresh.err)

	m, _ = update(t, m, refresh)
	assert.Len(t, m.items, 2)
	assert.Equal(t, 1, m.cursor, "cursor must clamp after the last item is deleted")
}

func TestLoadingIgnoresInputButAllowsQuit(t *testing.T) {
	m := newTestModel(&fakeStore{items: testItems()}, nil)
	m.loading = true

	next, cmd := update(t, m, key("down"))
	assert.Equal(t, 0, next.cursor, "movement must be ignored while loading")
	assert.Nil(t, cmd)

	_, cmd = update(t, m, key("q"))
	require.NotNil(t, cmd, "quit must work while loading")
	assert.IsType(t, tea.QuitMsg{}, cmd())
}

func TestEmptyListNoOps(t *testing.T) {
	m := newTestModel(&fakeStore{}, nil)

	for _, k := range []string{"enter", "d"} {
		next, cmd := update(t, m, key(k))
		assert.Nil(t, cmd, "key %q must be a no-op on an empty list", k)
		assert.False(t, next.loading)
	}
}

func TestRepoErrorSurfacesWithoutLosingItems(t *testing.T) {
	store := &fakeStore{items: testItems(), updateErr: errors.New("boom")}
	m := newTestModel(store, nil)

	m, cmd := update(t, m, key("enter"))
	require.NotNil(t, cmd)

	msg := cmd()
	refresh, ok := msg.(refreshMsg)
	require.True(t, ok)
	require.Error(t, refresh.err)

	m, _ = update(t, m, refresh)
	assert.False(t, m.loading)
	assert.Error(t, m.err)
	assert.Len(t, m.items, 3, "items must survive a failed operation")
}

func TestFilterPersistsAcrossRefresh(t *testing.T) {
	store := &fakeStore{items: testItems()}
	notDone := func(it todo.Todo) bool { return !it.Done }
	m := newTestModel(store, notDone)

	// Toggle the item under the cursor; the refresh must drop it,
	// because the filter only keeps not-done items.
	m, cmd := update(t, m, key("enter"))
	require.NotNil(t, cmd)

	msg := cmd()
	refresh, ok := msg.(refreshMsg)
	require.True(t, ok)
	require.NoError(t, refresh.err)

	m, _ = update(t, m, refresh)
	assert.Len(t, m.items, 2, "filter must be re-applied on refresh")
	for _, it := range m.items {
		assert.False(t, it.Done)
	}
}

func TestViewRendersStatusAndHelp(t *testing.T) {
	store := &fakeStore{items: []todo.Todo{
		{ID: 1, Text: "Open task", Position: 1},
		{ID: 2, Text: "Closed task", Position: 2, Done: true},
	}}
	m := newTestModel(store, nil)

	out := m.View()
	assert.Contains(t, out, "[ ]")
	assert.Contains(t, out, "[x]")
	assert.Contains(t, out, "Open task")
	assert.Contains(t, out, "q quit")
}

func TestViewEmptyList(t *testing.T) {
	m := newTestModel(&fakeStore{}, nil)
	assert.Contains(t, m.View(), "No tasks.")
}
