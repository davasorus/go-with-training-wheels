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

// Hellos returns a map that associates each of the named people with a greeting message.
func Hellos(names []string) (map[string]string, error) {
	//A map to associate names with messages.
	messages := make(map[string]string)
	//loop through the recieved slice of names, calling the hello function to get a message for each name.
	for _, name := range names {
		message, err := Hello(name)
		if err != nil {
			return nil, err // Return the error to the caller if a name is empty.
		}
		messages[name] = message // In the map, associate the retrieved message with the name.
	}
	return messages, nil
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
