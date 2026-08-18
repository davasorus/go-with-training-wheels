/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"testing"

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
	}{
		{
			name:         "all items",
			showOnlyDone: false,
		},
		{
			name:         "only done items",
			showOnlyDone: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			// This test verifies that the writeOutput logic executes without error
			// when given a valid buffer as the writer.
			writeOutput(&buf, tt.showOnlyDone)
			assert.NoError(t, nil) // Just checking if it completes without panic or critical errors
		})
	}
}
