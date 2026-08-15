package todo

import (
	"testing"
)

func TestRepository(t *testing.T) {
	// NewRepository calls Database.NewStore(), which handles initialization if needed,
	// but since we want to remove direct usage of Database.DB, we just call the constructor.
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Test SaveItems and ListItems integration
	testItems := []Todo{
		{Text: "Test Item 1", Priority: 2, Position: 1},
		{Text: "Test Item 2", Priority: 1, Position: 2},
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

	// Test UpdateItemStatus on valid item
	err = repo.UpdateItemStatus(1, true) // Update the first item
	if err != nil {
		t.Fatalf("Failed to update item status: %v", err)
	}

	// Verify that only the first item's state was changed in the database
	itemsAfterUpdate, _ := repo.ListItems()
	if !itemsAfterUpdate[0].Done {
		t.Errorf("Expected item 1 to be marked as done")
	}
	if itemsAfterUpdate[1].Done {
		t.Errorf("Unexpected state: Item 2 was also marked as done")
	}

	// Test UpdateItemStatus of a non-existent item
	err = repo.UpdateItemStatus(999, true)
	if err == nil {
		t.Errorf("Expected an error when updating a non-existent item (id 999), but got none")
	}

	// Test SaveItems with duplicates (Upsert check)
	testItemsDuplicate := []Todo{
		{Text: "Test Item 1", Priority: 2, Position: 1}, // Already exists
		{Text: "New Item", Priority: 0, Position: 3},
	}
	err = repo.SaveItems(testItemsDuplicate)
	if err != nil {
		t.Fatalf("Failed to save items with duplicates (upsert): %v", err)
	}

	itemsAfterUpsert, _ := repo.ListItems()
	if len(itemsAfterUpsert) != 2 {
		t.Errorf("Expected 2 unique items after upsert, got %d", len(itemsAfterUpsert))
	}
}
