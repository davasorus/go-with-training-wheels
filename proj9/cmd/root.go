/*
Copyright © 2026 Sean Davitt
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/davasorus/tri/Database"
	"github.com/spf13/cobra"
)

var dataFile string

// rootCmd represents the base command of the tri tool.
var rootCmd = &cobra.Command{
	Use:   "tri",
	Short: "A brief description of your application",
	Long: `A longer description of the application.
Example:
tri add -p high task name`,
}

// Execute initializes the database and runs the command tree.
func Execute() {
	if err := Database.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Database connection error: %v\n", err)
		os.Exit(1)
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Define global flags here.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
