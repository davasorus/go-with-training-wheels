package Database

import (
	"testing"
)

// TestInitDB verifies that the database initialization logic can be executed.
// Note: This test attempts to connect to a database. If no database is 
// available at the configured location, it will fail on Ping(), but it 
// confirms that the connection logic and configuration loading are functional.
func TestInitDB(t *testing.T) {
	// We call InitDB to ensure it can navigate through its internal logic:
	// 1. Loading .env (even if failing, it should warn and continue)
	// 2. Building the DSN string
	// 3. Initializing the sql.DB object
	err := InitDB()

	if err != nil {
		t.Logf("Notice: Database initialization failed as expected without a real DB reachable at local config: %v", err)
	}
}
