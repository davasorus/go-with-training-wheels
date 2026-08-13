package todo

import (
	"testing"

	"github.com/davasorus/tri/Database"
)

func TestRepository(t *testing.T) {
	// Ensure DB is initialized before running tests
	if Database.DB == nil {
		err := Database.InitDB()
		if err != nil {
			t.Fatalf("Failed to initialize database: %v", err)
		}
	}

	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Test SaveItems and ListItems integration
	testItems := []Todo{
		{Text: "Test Item 1", Priority: 2, position: 1},
		{Text: "Test Item 2", Priority: 1, position: 2},
	}
	err = repo.SaveItems(testItems)
	if err != nil {
		t.Fatalf("Failed to save items: %v", err)
	}

	items, err := repo.ListItems()
	if err != nil {
		t.Fatalf("Failed to list items: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Verify data content matches what we inserted
	if items[0].Text != "Test Item 1" || items[0].Priority != 2 || items[0].Done != false {
		t.Errorf("Item 1 mismatch: %+v", items[0])
	}
	if items[1].Text != "Test Item 2" || items[1].Priority != 1 || items[1].Done != false {
		t.Errorf("Item 2 mismatch: %+v", items[1])
	}

	// Test UpdateItemStatus
	err = repo.UpdateItemStatus(1, true) // Update the first item
	if err != nil {
		t.Fatalf("Failed to update item status: %v", err)
	}

	// Verify that only the first item's state was changed in the database
	itemsAfterUpdate, err := repo.ListItems()
	if err != nil {
		t.Fatalf("Failed to list items after update: %v", err)
	}

	if !itemsAfterUpdate[0].Done {
		t.Errorf("Expected item 1 to be marked as done")
	}
	if itemsAfterUpdate[1].Done {
		t.Errorf("Unexpected state: Item 2 was also marked as done")
	}
}
