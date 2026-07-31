package greetings

import "fmt"

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	//:= is a short variable declaration, which declares and initializes the variable in one line.
	//can also be broken up into multiple lines for readability
	//var message string
	//message = fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
