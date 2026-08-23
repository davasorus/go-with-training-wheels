package cmd

import (
	"fmt"
	"os"

	"github.com/davasorus/tri/Database"
	"github.com/davasorus/tri/todo"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command of the tri tool.
var rootCmd = &cobra.Command{
	Use:          "tri",
	SilenceUsage: true,
	Short:        "A brief description of your application",
	Long: `A longer description of the application.
Example:
tri add -p high task name`,
}

// Execute initializes the database and runs the command tree.
func Execute() {
	store, err := Database.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database connection error: %v\n", err)
		os.Exit(1)
	}

	adapter := todo.NewTodoStoreAdapter(store)
	repo := todo.NewRepository(adapter)

	// Register subcommands with the repository injected.
	rootCmd.AddCommand(addCmd(repo))
	rootCmd.AddCommand(listCmd(repo))
	rootCmd.AddCommand(doneCmd(repo))
	rootCmd.AddCommand(undoCmd(repo))
	rootCmd.AddCommand(editCmd(repo))
	rootCmd.AddCommand(moveCmd(repo))
	rootCmd.AddCommand(exportCmd(repo))
	rootCmd.AddCommand(clearCmd(repo))
	rootCmd.AddCommand(deleteCmd(repo))
	rootCmd.AddCommand(NewMigrateDbCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
