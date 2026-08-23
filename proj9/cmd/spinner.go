package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// spinnerDelay is how long an action may run before a spinner appears.
// Fast operations finish silently instead of flashing one frame.
const spinnerDelay = 150 * time.Millisecond

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

// WaitSpinner runs action, showing a spinner only if it takes longer
// than spinnerDelay. It leaves the line clean; callers print their own
// result message.
func WaitSpinner(message string, action func() error) error {
	done := make(chan error, 1)
	go func() { done <- action() }()

	select {
	case err := <-done:
		return err
	case <-time.After(spinnerDelay):
	}

	stop := make(chan struct{})
	go Spinner(message, stop)
	err := <-done
	close(stop)

	// Clear the spinner line.
	fmt.Printf("\r%s\r", strings.Repeat(" ", len(message)+2))
	return err
}

// confirm asks the user a yes/no question on the terminal. It returns
// true only for an explicit yes.
func confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}
