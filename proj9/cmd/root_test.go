/*
Copyright © 2026 Sean Davitt
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	// Verify that the root command is correctly initialized.
	assert.NotNil(t, rootCmd)
}
