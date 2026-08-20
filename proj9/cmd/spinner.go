package cmd

import (
	"fmt"
	"time"
)

// Spinner prints a spinning character to the terminal.
func Spinner(message string) {
	spinner := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		fmt.Printf("\r%s %s", message, spinner[i%4])
		time.Sleep(100 * time.Millisecond)
		i++
	}
}

// WaitSpinner is a helper to show progress while an action is performed.
func WaitSpinner(message string, action func() error) error {
	fmt.Printf("%s... ", message)

	// Run spinner in background if it's a potentially long-running task.
	go Spinner(message)

	err := action()

	// Clear the line after completion to remove the spinner and show status.
	fmt.Printf("\r%s %s\n", message, func() string {
		if err != nil {
			return "failed"
		}
		return "done"
	}())

	return err
}
