package main

import (
	"fmt"
	"log"

	"github.com/davasorus/greetings" // Update with your actual module path
)

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing the time, source file, and line number.
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"Gladys", "Samantha", "Darrin"}

	// Request a greeting
	messages, err := greetings.Hellos(names)
	// If an error was returned, print it to the console and exit the program.
	if err != nil {
		log.Fatal(err)
	}

	//Otherwise print
	fmt.Println(messages)
}
