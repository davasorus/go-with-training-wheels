package Database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Config holds database connection parameters.
type Config struct {

	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// DB is the global database connection pool.
// DB is the global connection pool for the database.
var DB *sql.DB

// InitDB initializes the database connection.
func InitDB() error {
	// Load environment variables from the specified path
	err := godotenv.Load("Database/.env")
	if err != nil {
		fmt.Println("Warning: could not load .env file, using system environment variables.")
	}

	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	if cfg.DBName == "" {
		cfg.DBName = "todos"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	var err2 error
	DB, err2 = sql.Open("postgres", dsn)
	if err2 != nil {
		return fmt.Errorf("failed to open database: %w", err2)
	}

	if err = DB.Ping(); err != nil {
		// If ping fails, the database may not exist. Try to create it if missing.
		fmt.Printf("Ping failed. Attempt to ensure database %s exists...\n", cfg.DBName)

		// Connect to the default 'postgres' database to check or create the target.
		adminDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password)

		tmpDB, errTmp := sql.Open("postgres", adminDSN)
		if errTmp != nil {
			return fmt.Errorf("failed to open administrative connection: %w", errTmp)
		}
		defer tmpDB.Close()

		// Check if the database exists
		var exists bool
		query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname='%s')`, cfg.DBName)
		err = tmpDB.QueryRow(query).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if database exists: %w", err)
		}

		if !exists {
			fmt.Printf("Database %s not found. Create it now...\n", cfg.DBName)
			// Note: You cannot use "CREATE DATABASE" in a transaction block or with other SQL features.
			// Execute it directly.
			_, err = tmpDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
			if err != nil {
				return fmt.Errorf("failed to create database: %w", err)
			}
		}

		// Now try to ping the original connection again
		if err = DB.Ping(); err != nil {
			return fmt.Errorf("failed to ping database after creation attempt: %w", err)
		}
	}

	return nil
}
