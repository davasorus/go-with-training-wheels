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

func TestPriorityFlagSet(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{"low", 0, false},
		{"LOW", 0, false},
		{"medium", 1, false},
		{"med", 1, false},
		{"high", 2, false},
		{"High", 2, false},
		{"0", 0, false},
		{"1", 1, false},
		{"2", 2, false},
		{"3", 0, true},
		{"urgent", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var p priorityFlag
			err := p.Set(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, int(p))
		})
	}
}

func TestParseDueDate(t *testing.T) {
	got, err := parseDueDate("")
	assert.NoError(t, err)
	assert.Nil(t, got)

	got, err = parseDueDate("today")
	assert.NoError(t, err)
	assert.NotNil(t, got)

	got, err = parseDueDate("Tomorrow")
	assert.NoError(t, err)
	assert.NotNil(t, got)

	got, err = parseDueDate("2026-09-01")
	assert.NoError(t, err)
	if assert.NotNil(t, got) {
		assert.Equal(t, "2026-09-01", got.Format("2006-01-02"))
	}

	_, err = parseDueDate("next week")
	assert.Error(t, err)

	_, err = parseDueDate("09/01/2026")
	assert.Error(t, err)
}
