package cmd

import (
	"fmt"
	"time"
)

// Spinner prints a spinning character to the terminal until stop is closed.
func Spinner(message string, stop <-chan struct{}) {
	frames := []string{"|", "/", "-", "\\"}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			fmt.Printf("\r%s %s", message, frames[i%4])
			i++
		}
	}
}

// WaitSpinner is a helper to show progress while an action is performed.
func WaitSpinner(message string, action func() error) error {
	stop := make(chan struct{})
	go Spinner(message, stop)

	err := action()
	close(stop)

	status := "done"
	if err != nil {
		status = "failed"
	}
	fmt.Printf("\r%s... %s\n", message, status)

	return err
}
