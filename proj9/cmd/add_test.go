/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"testing"

	"github.com/davasorus/tri/todo"
	"github.com/stretchr/testify/assert"
)

func TestPrepareAddItems(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		priority int
		expected []todo.Todo
	}{
		{
			name:     "single item",
			args:     []string{"Task1"},
			priority: 2,
			expected: []todo.Todo{
				{Text: "Task1", Priority: 2, Position: 1},
			},
		},
		{
			name:     "multiple items",
			args:     []string{"Task1", "Task2"},
			priority: 1,
			expected: []todo.Todo{
				{Text: "Task1", Priority: 1, Position: 1},
				{Text: "Task2", Priority: 1, Position: 2},
			},
		},
		{
			name:     "empty string input",
			args:     []string{""},
			priority: 0,
			expected: []todo.Todo{
				{Text: "", Priority: 0, Position: 1},
			},
		},
		{
			name:     "empty items list",
			args:     []string{},
			priority: 0,
			expected: []todo.Todo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prepareAddItems(tt.args, tt.priority, nil)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAddCmd(t *testing.T) {
	// Verify that the command is registered.
	assert.NotNil(t, addCmd)
}
