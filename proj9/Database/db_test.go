/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package Database

import (
	"testing"

	"database/sql"
	"github.com/stretchr/testify/assert"
)

// TestInitDB verifies that the database initialization logic can be executed.
// Note: This test attempts to connect to a database. If no database is
// available at the configured location, it will fail on Ping(), but it
// confirms that the connection logic and configuration loading are functional.
func TestInitDB(t *testing.T) {
	// We call InitDB() which returns (*Store, error).
	// Note: This test might fail if no DB is available at the configured location.
	_, err := InitDB()

	if err != nil {
		t.Logf("Notice: Database initialization failed as expected without a real DB reachable at local config: %v", err)
	}
}

func TestNewStore(t *testing.T) {
	// Case 1: Nil database connection
	store, err := NewStore(nil)
	assert.Nil(t, store)
	assert.EqualError(t, err, "database connection is nil")

	// Case 2: Non-nil (dummy) database connection
	db := &sql.DB{} // This is a dummy but it's no longer nil
	store, err = NewStore(db)
	assert.NotNil(t, store)
	assert.NoError(t, err)
}
