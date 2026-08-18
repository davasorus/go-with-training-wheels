package todo

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// Todo represents a single task in the list.
type Todo struct {
	Text     string
	Priority int
	Position int
	Done     bool
	DueDate  *time.Time // Use a pointer to handle NULL values from the database
}

// ByPri is a type for sorting tasks by priority and position.
type ByPri []Todo

func (a ByPri) Len() int      { return len(a) }
func (a ByPri) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPri) Less(i, j int) bool {
	if a[i].Done && !a[j].Done {
		return a[i].Done
	}
	if a[i].Priority == a[j].Priority {
		return a[i].Position < a[j].Position
	}
	return a[i].Position < a[j].Position
}

// SaveItems saves a slice of Todo items to the database using an UPSERT pattern.
func SaveItems(r TodoStore, items []Todo) error {
	err := r.SaveItems(items)
	if err != nil {
		return fmt.Errorf("failed to save items: %w", err)
	}
	slog.Debug("Saved items to database")
	return nil
}

// LoadItems retrieves all Todo items from the database.
func LoadItems(r TodoStore) ([]Todo, error) {
	items, err := r.ListItems()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	slog.Debug("Loaded items from database")
	return items, nil
}

func (i *Todo) SetPriority(pri int) {
	switch pri {
	case 0:
		i.Priority = 0
	case 1:
		i.Priority = 1
	case 2:
		i.Priority = 2
	default:
		i.Priority = 0
	}
}

// PrettyP returns the priority level as a string.
func (i *Todo) PrettyP() string {
	switch i.Priority {
	case 0:
		return "Low"
	case 1:
		return "Medium"
	case 2:
		return "High"
	default:
		return "Low"
	}
}

// Label returns the position of the task.
func (i *Todo) Label() string {
	return strconv.Itoa(i.Position) + "."
}

// DoneStatus returns the status of the task.
func (i *Todo) DoneStatus() string {
	if i.Done {
		return "Done"
	}
	return "Not Done"
}

// String returns a formatted string representation of the task.
func (i *Todo) String() string {
	return fmt.Sprintf("%s %s [%s] %s", i.Label(), i.Text, i.PrettyP(), i.DoneStatus())
}
