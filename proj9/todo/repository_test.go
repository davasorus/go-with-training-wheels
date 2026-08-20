/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package todo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockStore is used for unit testing Repository logic without requiring a database connection.
type mockStore struct {
	items       map[int]Todo
	saveErr     error
	updateErr   error
	notFoundErr error
}

// ListItems returns all items from the mock store.
func (m *mockStore) ListItems() ([]Todo, error) {
	var res []Todo
	for _, item := range m.items {
		res = append(res, item)
	}
	return res, nil
}

// SaveItems saves items to the mock store. It simulates a save operation and can return a predefined error.
func (m *mockStore) SaveItems(items []Todo) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	for _, item := range items {
		m.items[item.Position] = item
	}
	return nil
}

// UpdateItemStatus updates the status of a specific todo in the mock store. It can simulate errors for testing.
func (m *mockStore) UpdateItemStatus(id int, done bool) error {
	if m.notFoundErr != nil && id == 999 {
		return m.notFoundErr
	}
	if m.updateErr != nil {
		return m.updateErr
	}
	if item, ok := m.items[id]; ok {
		item.Done = done
		m.items[id] = item
	}
	return nil
}

// DeleteItem removes an item from the mock store.
func (m *mockStore) DeleteItem(id int) error {
	delete(m.items, id)
	return nil
}

func TestRepository_Mocked(t *testing.T) {
	mockData := map[int]Todo{
		1: {Text: "Test Item 1", Priority: 2, Position: 1, Done: false},
		2: {Text: "Test Item 2", Priority: 1, Position: 2, Done: false},
	}

	m := &mockStore{
		items:       mockData,
		notFoundErr: errors.New("not found"),
	}

	repo := &Repository{store: m}

	t.Run("ListItems", func(t *testing.T) {
		items, err := repo.ListItems()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(items))
	})

	t.Run("ListItems_Empty", func(t *testing.T) {
		mEmpty := &mockStore{items: make(map[int]Todo)}
		repoEmpty := &Repository{store: mEmpty}
		items, err := repoEmpty.ListItems()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(items))
	})

	t.Run("SaveItems_Success", func(t *testing.T) {
		testData := []Todo{
			{Text: "New Item", Priority: 3, Position: 3, Done: false},
		}
		err := repo.SaveItems(testData)
		assert.NoError(t, err)
	})

	t.Run("SaveItems_Failure", func(t *testing.T) {
		mFail := &mockStore{saveErr: errors.New("database error")}
		repoFail := &Repository{store: mFail}
		err := repoFail.SaveItems([]Todo{{Text: "fail"}})
		assert.Error(t, err)
	})

	t.Run("UpdateItemStatus_Success", func(t *testing.T) {
		err := repo.UpdateItemStatus(1, true)
		assert.NoError(t, err)
		assert.True(t, m.items[1].Done)
	})

	t.Run("UpdateItemStatus_NotFound", func(t *testing.T) {
		err := repo.UpdateItemStatus(999, true)
		assert.Error(t, err)
	})

	t.Run("UpdateItemStatus_InternalError", func(t *testing.T) {
		mErr := &mockStore{updateErr: errors.New("internal error")}
		repoErr := &Repository{store: mErr}
		err := repoErr.UpdateItemStatus(1, true)
		assert.Error(t, err)
	})
}
