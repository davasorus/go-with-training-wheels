package models

import (
	"testing"
)

func TestModelLoading(t *testing.T) {
	// This test ensures that the model is correctly defined and accessible.
	m := Todo{
		Text:     "Test Task",
		Priority: 1,
		Position: 10,
		Done:     false,
	}

	if m.Text != "Test Task" {
		t.Errorf("expected text 'Test Task', got '%s'", m.Text)
	}

	if m.Priority != 1 {
		t.Errorf("expected priority 1, got %d", m.Priority)
	}

	if m.Position != 10 {
		t.Errorf("expected position 10, got %d", m.Position)
	}

	if m.Done {
		t.Error("expected Done to be false")
	}
}
