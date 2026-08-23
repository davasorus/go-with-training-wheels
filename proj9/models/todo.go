package models

import (
	"time"
)

// Todo represents a task in the to-do list application.
type Todo struct {
	ID          int
	Text        string
	Priority    int
	Position    int
	Done        bool
	DueDate     *time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
	Tags        []string
}
