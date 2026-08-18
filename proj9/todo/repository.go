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
func NewRepository(store TodoStore) *Repository {
	return &Repository{store: store}
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

// TodoStoreAdapter wraps Database.Store and implements TodoStore by mapping models.Todo to todo.Todo.
type TodoStoreAdapter struct {
	*Database.Store
}

// NewTodoStoreAdapter creates a new instance of TodoStoreAdapter.
func NewTodoStoreAdapter(store *Database.Store) *TodoStoreAdapter {
	return &TodoStoreAdapter{Store: store}
}

// ListItems fetches all items from the database and converts them to the todo.Todo type.
func (a *TodoStoreAdapter) ListItems() ([]Todo, error) {
	items, err := a.Store.ListItems()
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
			DueDate:  item.DueDate,
		}
	}
	return res, nil
}

// SaveItems converts todo.Todo items to models.Todo and saves them to the database.
func (a *TodoStoreAdapter) SaveItems(items []Todo) error {
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

// UpdateItemStatus updates the status of a specific todo in the database.
func (a *TodoStoreAdapter) UpdateItemStatus(id int, done bool) error {
	return a.Store.UpdateItemStatus(id, done)
}
