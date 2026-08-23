package ui

import (
	"errors"
	"fmt"
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

func (f *fakeStore) SaveItems(items []todo.Todo) error {
	for _, it := range items {
		it.ID = len(f.items) + 100
		f.items = append(f.items, it)
	}
	return nil
}

func (f *fakeStore) UpdateItem(item todo.Todo) error {
	for i := range f.items {
		if f.items[i].ID == item.ID {
			f.items[i] = item
			return nil
		}
	}
	return errors.New("not found")
}

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

func TestAddMode(t *testing.T) {
	store := &fakeStore{items: testItems()}
	m := newTestModel(store, nil)

	// 'a' enters input mode.
	m, _ = update(t, m, key("a"))
	assert.Equal(t, modeAdd, m.mode)

	// Typed runes go to the input, not the list.
	m, _ = update(t, m, key("x"))
	assert.Equal(t, "x", m.input.Value())
	assert.Equal(t, 0, m.cursor, "list keys must not fire in input mode")

	// Enter saves, exits input mode, and triggers a refresh.
	m, cmd := update(t, m, key("enter"))
	assert.Equal(t, modeNormal, m.mode)
	assert.True(t, m.loading)
	require.NotNil(t, cmd)

	msg := cmd()
	refresh, ok := msg.(refreshMsg)
	require.True(t, ok)
	require.NoError(t, refresh.err)

	m, _ = update(t, m, refresh)
	assert.Len(t, m.items, 4, "new task must appear after refresh")
}

func TestAddModeEscCancels(t *testing.T) {
	m := newTestModel(&fakeStore{items: testItems()}, nil)

	m, _ = update(t, m, key("a"))
	m, _ = update(t, m, key("x"))
	m, cmd := update(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	assert.Equal(t, modeNormal, m.mode)
	assert.Nil(t, cmd)
	assert.Empty(t, m.input.Value(), "esc must clear the input")
	assert.Len(t, m.items, 3)
}

func TestAddModeEmptyInputIsNoOp(t *testing.T) {
	m := newTestModel(&fakeStore{items: testItems()}, nil)

	m, _ = update(t, m, key("a"))
	m, cmd := update(t, m, key("enter"))
	assert.Equal(t, modeNormal, m.mode)
	assert.False(t, m.loading)
	assert.Nil(t, cmd, "saving empty text must be a no-op")
}

func TestCursorFollowsToggledItem(t *testing.T) {
	// Toggling done re-sorts the list; the cursor must follow the item.
	store := &fakeStore{items: testItems()}
	m := newTestModel(store, nil)
	m.items = append([]todo.Todo{}, store.items...) // ID 10, 20, 30

	// Cursor on the last item (ID 30). Toggling it done moves it to the
	// top of the sort (done items sort first).
	m.cursor = 2
	m, cmd := update(t, m, key("enter"))
	require.NotNil(t, cmd)
	refresh := cmd().(refreshMsg)
	require.NoError(t, refresh.err)

	m, _ = update(t, m, refresh)
	assert.Equal(t, 30, m.items[m.cursor].ID, "cursor must stay on the toggled item after re-sort")
	assert.NotEqual(t, 2, m.cursor, "the item must actually have moved")
}

func TestEditMode(t *testing.T) {
	store := &fakeStore{items: testItems()}
	m := newTestModel(store, nil)

	m, _ = update(t, m, key("e"))
	assert.Equal(t, modeEdit, m.mode)
	assert.Equal(t, "Task A", m.input.Value(), "edit must prefill the current text")

	// Append text and save.
	m, _ = update(t, m, key("!"))
	m, cmd := update(t, m, key("enter"))
	assert.Equal(t, modeNormal, m.mode)
	require.NotNil(t, cmd)

	refresh := cmd().(refreshMsg)
	require.NoError(t, refresh.err)
	m, _ = update(t, m, refresh)

	var found bool
	for _, it := range m.items {
		if it.ID == 10 {
			found = true
			assert.Equal(t, "Task A!", it.Text)
		}
	}
	assert.True(t, found)
}

func TestPriorityKeys(t *testing.T) {
	store := &fakeStore{items: testItems()}
	m := newTestModel(store, nil)

	// Cursor on ID 10 (priority 2). Press "1" -> low.
	m, cmd := update(t, m, key("1"))
	assert.True(t, m.loading)
	require.NotNil(t, cmd)

	refresh := cmd().(refreshMsg)
	require.NoError(t, refresh.err)
	m, _ = update(t, m, refresh)

	for _, it := range m.items {
		if it.ID == 10 {
			assert.Equal(t, 0, it.Priority)
		}
	}
	assert.Equal(t, 10, m.items[m.cursor].ID, "cursor must follow the reprioritized item")
}

func TestScrollingWindow(t *testing.T) {
	items := make([]todo.Todo, 30)
	for i := range items {
		items[i] = todo.Todo{ID: i + 1, Text: fmt.Sprintf("Task %02d", i+1), Position: i + 1}
	}
	m := newTestModel(&fakeStore{items: items}, nil)
	m.items = items

	// Simulate a small terminal: 10 lines tall -> 3 visible rows.
	m, _ = update(t, m, tea.WindowSizeMsg{Width: 80, Height: 10})

	// Walk the cursor beyond the window; the offset must follow.
	for i := 0; i < 10; i++ {
		m, _ = update(t, m, key("down"))
	}
	assert.Equal(t, 10, m.cursor)
	assert.Greater(t, m.offset, 0, "window must scroll down with the cursor")

	out := m.View()
	assert.Contains(t, out, "more", "scroll indicators must render")
	assert.Contains(t, out, "Task 11", "cursor row must be visible")
	assert.NotContains(t, out, "Task 01", "scrolled-off rows must not render")

	// Walk back up; offset must follow.
	for i := 0; i < 10; i++ {
		m, _ = update(t, m, key("up"))
	}
	assert.Equal(t, 0, m.cursor)
	assert.Equal(t, 0, m.offset)
}
