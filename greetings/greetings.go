package greetings

import (
	"fmt"
	"math/rand"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
	//If no name was given, return an error with a message.
	if name == "" {
		return "", fmt.Errorf("empty name")
	}

	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf(randomFormat(), name)
	//:= is a short variable declaration, which declares and initializes the variable in one line.
	//can also be broken up into multiple lines for readability
	//var message string
	//message = fmt.Sprintf("Hi, %v. Welcome!", name)
	return message, nil
}

func randomFormat() string {
	// A slice of message formats.
	formats := []string{
		"Hi, %v. Welcome!",
		"Hello, %v! How are you?",
		"Great to see you, %v!",
		"Hey, %v! Nice to meet you!",
	}
	// Return a randomly selected format.
	return formats[rand.Intn(len(formats))]
}
