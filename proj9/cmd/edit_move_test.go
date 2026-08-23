package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davasorus/tri/todo"
)

// cmdFakeStore is an in-memory TodoStore for command tests.
type cmdFakeStore struct {
	items []todo.Todo
}

func (f *cmdFakeStore) ListItems() ([]todo.Todo, error) {
	out := make([]todo.Todo, len(f.items))
	copy(out, f.items)
	return out, nil
}

func (f *cmdFakeStore) SaveItems(items []todo.Todo) error { return nil }

func (f *cmdFakeStore) UpdateItem(item todo.Todo) error {
	for i := range f.items {
		if f.items[i].ID == item.ID {
			f.items[i] = item
			return nil
		}
	}
	return errors.New("not found")
}

func (f *cmdFakeStore) UpdateItemStatus(id int, done bool) error {
	for i := range f.items {
		if f.items[i].ID == id {
			f.items[i].Done = done
			return nil
		}
	}
	return errors.New("not found")
}

func (f *cmdFakeStore) DeleteItem(id int) error { return nil }

func (f *cmdFakeStore) byID(id int) todo.Todo {
	for _, it := range f.items {
		if it.ID == id {
			return it
		}
	}
	return todo.Todo{}
}

func editStore() *cmdFakeStore {
	return &cmdFakeStore{items: []todo.Todo{
		{ID: 1, Text: "High task", Priority: 2, Position: 1},
		{ID: 2, Text: "Low task", Priority: 0, Position: 2},
	}}
}

func TestExecuteEditText(t *testing.T) {
	store := editStore()
	// Sorted order: High task is row 1.
	err := executeEdit(store, "1", false, 0, "Renamed", "", nil)
	require.NoError(t, err)
	assert.Equal(t, "Renamed", store.byID(1).Text)
}

func TestExecuteEditPriorityAndDue(t *testing.T) {
	store := editStore()
	err := executeEdit(store, "2", true, 2, "", "2026-09-01", nil)
	require.NoError(t, err)
	it := store.byID(2)
	assert.Equal(t, 2, it.Priority)
	require.NotNil(t, it.DueDate)
	assert.Equal(t, "2026-09-01", it.DueDate.Format("2006-01-02"))
}

func TestExecuteEditClearDueAndTags(t *testing.T) {
	now := time.Now()
	store := &cmdFakeStore{items: []todo.Todo{
		{ID: 1, Text: "Task", Priority: 1, Position: 1, DueDate: &now, Tags: []string{"work"}},
	}}
	err := executeEdit(store, "1", false, 0, "", "none", []string{"none"})
	require.NoError(t, err)
	it := store.byID(1)
	assert.Nil(t, it.DueDate)
	assert.Empty(t, it.Tags)
}

func TestExecuteEditNothingToChange(t *testing.T) {
	err := executeEdit(editStore(), "1", false, 0, "", "", nil)
	assert.Error(t, err)
}

func TestExecuteMoveSwapsWithinPriority(t *testing.T) {
	store := &cmdFakeStore{items: []todo.Todo{
		{ID: 1, Text: "A", Priority: 1, Position: 1},
		{ID: 2, Text: "B", Priority: 1, Position: 2},
		{ID: 3, Text: "C", Priority: 1, Position: 3},
	}}
	// Move row 3 (C) up: C and B swap.
	err := executeMove(store, "3", "up")
	require.NoError(t, err)
	assert.Less(t, store.byID(3).Position, store.byID(2).Position)
}

func TestExecuteMoveRejectsBoundaries(t *testing.T) {
	store := &cmdFakeStore{items: []todo.Todo{
		{ID: 1, Text: "High", Priority: 2, Position: 1},
		{ID: 2, Text: "Low", Priority: 0, Position: 2},
	}}
	// Row 1 is already at the top.
	assert.Error(t, executeMove(store, "1", "up"))
	// Row 2 moving up would cross a priority boundary.
	assert.Error(t, executeMove(store, "2", "up"))
	// Bad direction.
	assert.Error(t, executeMove(store, "1", "sideways"))
}
