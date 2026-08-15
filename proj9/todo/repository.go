package todo

import (
	"fmt"

	"github.com/davasorus/tri/Database"
	"github.com/davasorus/tri/models"
)

type TodoStore interface {
	ListItems() ([]Todo, error)
	SaveItems(items []Todo) error
	UpdateItemStatus(id int, done bool) error
}

// Repository manages access to todo records in the database.
type Repository struct {
	store TodoStore
}

// NewRepository creates a new instance of the Repository.
func NewRepository() (*Repository, error) {
	dbStore, err := Database.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database store: %w", err)
	}
	return &Repository{store: &todoStoreAdapter{dbStore}}, nil
}

// todoStoreAdapter wraps Database.Store and implements TodoStore by mapping models.Todo to todo.Todo.
type todoStoreAdapter struct {
	*Database.Store
}

func (a *todoStoreAdapter) ListItems() ([]Todo, error) {
	items, err := a.Store.ListItems() // Returns []models.Todo
	if err != nil {
		return nil, err
	}

	res := make([]Todo, len(items))
	for i, item := range items {
		res[i] = Todo{
			Text:     item.Text,
			Priority: item.Priority,
			Position: item.Position,
			Done:     item.Done,
			DueDate:  item.DueDate, // Note: This assumes the underlying type of item.DueDate matches the target's field.
			// Wait, I need to check why this is failing.
		}
	}
	return res, nil
}

func (a *todoStoreAdapter) SaveItems(items []Todo) error {
	mItems := make([]models.Todo, len(items))
	for i, item := range items {
		mItems[i] = models.Todo{
			Text:     item.Text,
			Priority: item.Priority,
			Position: item.Position,
			Done:     item.Done,
			DueDate:  item.DueDate,
		}
	}
	return a.Store.SaveItems(mItems)
}

func (a *todoStoreAdapter) UpdateItemStatus(id int, done bool) error {
	return a.Store.UpdateItemStatus(id, done)
}

// ListItems fetches all items from the database.
func (r *Repository) ListItems() ([]Todo, error) {
	items, err := r.store.ListItems()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	return items, nil
}

// SaveItems persists a slice of todos into the database using a transaction.
func (r *Repository) SaveItems(items []Todo) error {
	err := r.store.SaveItems(items)
	if err != nil {
		return fmt.Errorf("failed to save items: %w", err)
	}
	return nil
}

// UpdateItemStatus updates the status of a specific todo by its ID.
func (r *Repository) UpdateItemStatus(id int, done bool) error {
	err := r.store.UpdateItemStatus(id, done)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	return nil
}
