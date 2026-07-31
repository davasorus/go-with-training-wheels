package main

import (
	"fmt"

	"github.com/davasorus/greetings" // Update with your actual module path
)

func main() {
	// Get a greeting message and print it.
	message := greetings.Hello("Gladys")
	fmt.Println(message)
}
