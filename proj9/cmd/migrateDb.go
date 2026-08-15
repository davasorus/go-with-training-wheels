package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// migrateDbCmd bootstraps the todos database (if missing) and applies
// all Liquibase changesets. Safe to run repeatedly.
var migrateDbCmd = &cobra.Command{
	Use:   "migrateDb",
	Short: "Create the todos database if needed and apply all Liquibase changesets",
	RunE: func(cmd *cobra.Command, args []string) error {
		// same config source as the rest of the app
		_ = godotenv.Load("Database/.env")

		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		steps := [][]string{
			{
				"update",
				"--changelog-file=bootstrap-changelog.yaml",
				fmt.Sprintf("--url=jdbc:postgresql://%s:%s/postgres", host, port),
			},
			{
				"update",
				"--changelog-file=changelog-root.yaml",
				fmt.Sprintf("--url=jdbc:postgresql://%s:%s/todos", host, port),
			},
		}

		for _, argv := range steps {
			c := exec.Command("liquibase", argv...)
			c.Dir = "Database/Changelogs"
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("liquibase %s (%s) failed: %w", argv[0], argv[1], err)
			}
		}

		fmt.Println("Database migration complete.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateDbCmd)
}
