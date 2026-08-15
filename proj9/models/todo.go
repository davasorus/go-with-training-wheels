package models

import (
	"time"
)

type Todo struct {
	ID       int
	Text     string
	Priority int
	Position int
	Done     bool
	DueDate  *time.Time
}
