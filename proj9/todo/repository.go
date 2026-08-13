package todo

import (
	"database/sql"
	"fmt"

	"github.com/davasorus/tri/Database"
)

// Repository defines the data access layer for Todo items.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new instance of the Todo repository.
func NewRepository() (*Repository, error) {
	if Database.DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	return &Repository{db: Database.DB}, nil
}

// ListItems fetches all items from the database.
func (r *Repository) ListItems() ([]Todo, error) {
	rows, err := r.db.Query(`SELECT id, text, priority, position, done FROM todos`)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []Todo
	for rows.Next() {
		var id int
		var item Todo
		err := rows.Scan(&id, &item.Text, &item.Priority, &item.position, &item.Done)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// SaveItems persists a slice of todos into the database using a transaction.
func (r *Repository) SaveItems(items []Todo) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.Exec(`INSERT INTO todos (text, priority, position, done) VALUES ($1, $2, $3, $4)`,
			item.Text, item.Priority, item.position, item.Done)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// UpdateItemStatus updates the status of a specific todo by its ID.
func (r *Repository) UpdateItemStatus(id int, done bool) error {
	res, err := r.db.Exec(`UPDATE todos SET done = $1 WHERE id = $2`, done, id)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no record found with id %d", id)
	}
	return nil
}
