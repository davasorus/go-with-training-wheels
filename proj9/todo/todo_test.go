package todo

import (
	"testing"
)

func TestTodoMethods(t *testing.T) {
	tests := []struct {
		name     string
		item     Todo
		expected string // For String() method
		priority int
		testFunc func(*testing.T, Todo)
	}{
		{
			name: "SetPriority Low",
			item: Todo{Text: "Test", Priority: 5},
			testFunc: func(t *testing.T, i Todo) {
				i.SetPriority(0)
				if i.Priority != 0 {
					t.Errorf("Expected priority 0, got %d", i.Priority)
				}
			},
		},
		{
			name: "SetPriority Medium",
			item: Todo{Text: "Test", Priority: 5},
			testFunc: func(t *testing.T, i Todo) {
				i.SetPriority(1)
				if i.Priority != 1 {
					t.Errorf("Expected priority 1, got %d", i.Priority)
				}
			},
		},
		{
			name: "SetPriority High",
			item: Todo{Text: "Test", Priority: 5},
			testFunc: func(t *testing.T, i Todo) {
				i.SetPriority(2)
				if i.Priority != 2 {
					t.Errorf("Expected priority 2, got %d", i.Priority)
				}
			},
		},
		{
			name: "PrettyP Logic",
			item: Todo{Text: "Test", Priority: 0},
			testFunc: func(t *testing.T, i Todo) {
				if i.PrettyP() != "Low" {
					t.Errorf("Expected Low, got %s", i.PrettyP())
				}
			},
		},
		{
			name: "PrettyP Medium",
			item: Todo{Text: "Test", Priority: 1},
			testFunc: func(t *testing.T, i Todo) {
				if i.PrettyP() != "Medium" {
					t.Errorf("Expected Medium, got %s", i.PrettyP())
				}
			},
		},
		{
			name: "PrettyP High",
			item: Todo{Text: "Test", Priority: 2},
			testFunc: func(t *testing.T, i Todo) {
				if i.PrettyP() != "High" {
					t.Errorf("Expected High, got %s", i.PrettyP())
				}
			},
		},
		{
			name: "Label Generation",
			item: Todo{Text: "Test", Position: 5},
			testFunc: func(t *testing.T, i Todo) {
				if i.Label() != "5." {
					t.Errorf("Expected 5., got %s", i.Label())
				}
			},
		},
		{
			name: "DoneStatus True",
			item: Todo{Text: "Test", Done: true},
			testFunc: func(t *testing.T, i Todo) {
				if i.DoneStatus() != "Done" {
					t.Errorf("Expected Done, got %s", i.DoneStatus())
				}
			},
		},
		{
			name: "DoneStatus False",
			item: Todo{Text: "Test", Done: false},
			testFunc: func(t *testing.T, i Todo) {
				if i.DoneStatus() != "Not Done" {
					t.Errorf("Expected Not Done, got %s", i.DoneStatus())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.item)
		})
	}
}

func TestToString(t *testing.T) {
	item := Todo{
		Text:     "Test Item",
		Priority: 2,
		Position: 1,
		Done:     true,
	}
	expected := "1. Test Item [High] Done"
	if item.String() != expected {
		t.Errorf("Expected %s, got %s", expected, item.String())
	}
}
