package Database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/davasorus/tri/models"
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

// Store wraps the sql.DB and provides methods for database operations.
type Store struct {
	db *sql.DB
}

// NewStore creates a new instance of the Store from an existing connection.
func NewStore(db *sql.DB) (*Store, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return &Store{db: db}, nil
}

// ListItems fetches all items from the database.
func (s *Store) ListItems() ([]models.Todo, error) {
	rows, err := s.db.Query(`SELECT id, text, priority, position, done, due_date FROM todos`)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []models.Todo
	for rows.Next() {
		var item models.Todo
		err := rows.Scan(&item.ID, &item.Text, &item.Priority, &item.Position, &item.Done, &item.DueDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while iterating items: %w", err)
	}

	return items, nil
}

// SaveItems persists a slice of todos into the database using a transaction.
func (s *Store) SaveItems(items []models.Todo) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	for _, item := range items {
		_, err := tx.Exec(`INSERT INTO todos (text, priority, position, done, due_date) 
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (text) DO UPDATE SET 
				priority = EXCLUDED.priority,
				position = EXCLUDED.position,
				done = EXCLUDED.done,
				due_date = EXCLUDED.due_date`,
			item.Text, item.Priority, item.Position, item.Done, item.DueDate)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// UpdateItemStatus updates the status of a specific todo by its ID.
func (s *Store) UpdateItemStatus(id int, done bool) error {
	res, err := s.db.Exec(`UPDATE todos SET done = $1 WHERE id = $2`, done, id)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no record found with id %d", id)
	}
	return nil
}

// InitDB initializes the database connection and returns a Store.
func InitDB() (*Store, error) {
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

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		// If ping fails, the database may not exist. Try to create it if missing.
		fmt.Printf("Ping failed. Attempt to ensure database %s exists...\n", cfg.DBName)

		// Connect to the default 'postgres' database to check or create the target.
		adminDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password)

		tmpDB, errTmp := sql.Open("postgres", adminDSN)
		if errTmp != nil {
			return nil, fmt.Errorf("failed to open administrative connection: %w", errTmp)
		}
		defer tmpDB.Close()

		// Check if the database exists
		var exists bool
		query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname='%s')`, cfg.DBName)
		err = tmpDB.QueryRow(query).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check if database exists: %w", err)
		}

		if !exists {
			fmt.Printf("Database %s not found. Create it now...\n", cfg.DBName)
			// Note: You cannot use "CREATE DATABASE" in a transaction block or with other SQL features.
			// Execute it directly.
			_, err = tmpDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
			if err != nil {
				return nil, fmt.Errorf("failed to create database: %w", err)
			}
		}

		// Now try to ping the original connection again
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping database after creation attempt: %w", err)
		}
	}

	return NewStore(db)
}
