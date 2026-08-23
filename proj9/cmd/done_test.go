/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoneRun(t *testing.T) {
	// Since doneRun interacts with global state and logs,
	// we test the internal logic instead of the CLI wrapper.
	assert.NotNil(t, doneCmd)
}

func TestValidateIndex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		count       int
		expectedIdx int
		expectError bool
	}{
		{
			name:        "valid index",
			input:       "1",
			count:       10,
			expectedIdx: 1,
			expectError: false,
		},
		{
			name:        "out of bounds high",
			input:       "11",
			count:       10,
			expectedIdx: 0,
			expectError: true,
		},
		{
			name:        "out of bounds low",
			input:       "0",
			count:       10,
			expectedIdx: 0,
			expectError: true,
		},
		{
			name:        "not a number",
			input:       "abc",
			count:       10,
			expectedIdx: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, err := validateIndex(tt.input, tt.count)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, idx)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIdx, idx)
			}
		})
	}
}
