/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"testing"

	"github.com/davasorus/tri/todo"
	"github.com/stretchr/testify/assert"
)

func TestListCmd(t *testing.T) {
	// Verify the command is registered.
	assert.NotNil(t, listCmd)
}

func TestWriteOutput(t *testing.T) {
	tests := []struct {
		name         string
		showOnlyDone bool
		items        []todo.Todo
		expected     []string // Substrings to look for in the output
	}{
		{
			name:         "all items",
			showOnlyDone: false,
			items: []todo.Todo{
				{Text: "Task 1", Priority: 1, Position: 1, Done: false},
				{Text: "Task 2", Priority: 2, Position: 2, Done: true},
			},
			expected: []string{"1.", "Task 1", "Not Done", "2.", "Task 2", "Done"},
		},
		{
			name:         "only done items (none done)",
			showOnlyDone: true,
			items: []todo.Todo{
				{Text: "Task 1", Priority: 1, Position: 1, Done: false},
			},
			expected: []string{},
		},
		{
			name:         "only done items (some done)",
			showOnlyDone: true,
			items: []todo.Todo{
				{Text: "Task 1", Priority: 1, Position: 1, Done: false},
				{Text: "Task 2", Priority: 2, Position: 2, Done: true},
			},
			expected: []string{"2.", "Task 2", "Done"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderList(&buf, tt.items, tt.showOnlyDone)
			output := buf.String()

			if len(tt.expected) == 0 {
				assert.Empty(t, output)
			} else {
				for _, exp := range tt.expected {
					assert.Contains(t, output, exp)
				}
			}
		})
	}
}
