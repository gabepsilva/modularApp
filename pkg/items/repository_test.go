package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()

	// Assert the repository is initialized correctly
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.items)
	assert.Equal(t, 1, repo.nextID)
	assert.Empty(t, repo.items)
}

func TestRepositoryCreateItem(t *testing.T) {
	repo := NewRepository()

	// Create an item
	item := repo.Create("Test Item", "This is a test item")

	// Assert the item was created with correct values
	assert.Equal(t, "1", item.ID)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, "This is a test item", item.Content)

	// Verify the item is stored in the repository
	assert.Len(t, repo.items, 1)
	assert.Contains(t, repo.items, "1")

	// Create another item to verify ID increment
	item2 := repo.Create("Another Item", "This is another test item")
	assert.Equal(t, "2", item2.ID)
	assert.Len(t, repo.items, 2)
}

func TestRepositoryGetAll(t *testing.T) {
	repo := NewRepository()

	// Repository should be empty initially
	items := repo.GetAll()
	assert.Empty(t, items)

	// Add a few items
	repo.Create("Item 1", "Content 1")
	repo.Create("Item 2", "Content 2")
	repo.Create("Item 3", "Content 3")

	// Get all items
	items = repo.GetAll()

	// Assert we have the correct number of items
	assert.Len(t, items, 3)

	// Verify item values - order is not guaranteed since we're using a map
	// so we'll collect all names and verify all three are present
	names := make([]string, 0, 3)
	for _, item := range items {
		names = append(names, item.Name)
	}

	assert.Contains(t, names, "Item 1")
	assert.Contains(t, names, "Item 2")
	assert.Contains(t, names, "Item 3")
}

func TestRepositoryGetByID(t *testing.T) {
	repo := NewRepository()

	// Add an item
	repo.Create("Test Item", "This is a test item")

	// Test getting an existing item
	item, err := repo.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", item.ID)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, "This is a test item", item.Content)

	// Test getting a non-existent item
	_, err = repo.GetByID("999")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)
}

func TestRepositoryUpdate(t *testing.T) {
	repo := NewRepository()

	// Add an item
	repo.Create("Original Name", "Original content")

	// Test updating an existing item's name
	updatedItem, err := repo.Update("1", "Updated Name", "")
	assert.NoError(t, err)
	assert.Equal(t, "1", updatedItem.ID)
	assert.Equal(t, "Updated Name", updatedItem.Name)
	assert.Equal(t, "Original content", updatedItem.Content)

	// Test updating an existing item's content
	updatedItem, err = repo.Update("1", "", "Updated content")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedItem.Name)
	assert.Equal(t, "Updated content", updatedItem.Content)

	// Test updating both name and content
	updatedItem, err = repo.Update("1", "New Name", "New content")
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updatedItem.Name)
	assert.Equal(t, "New content", updatedItem.Content)

	// Test updating a non-existent item
	_, err = repo.Update("999", "Test", "Test")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)
}

func TestRepositoryDelete(t *testing.T) {
	repo := NewRepository()

	// Add an item
	repo.Create("Test Item", "This is a test item")

	// Verify item exists
	assert.Len(t, repo.items, 1)

	// Delete the item
	err := repo.Delete("1")
	assert.NoError(t, err)

	// Verify item is deleted
	assert.Empty(t, repo.items)

	// Try to delete a non-existent item
	err = repo.Delete("1")
	assert.Error(t, err)
	assert.Equal(t, ErrItemNotFound, err)
}

func TestRepositoryConcurrentAccess(t *testing.T) {
	// This test ensures that the mutex protects against race conditions
	// Note: To fully test this, run with the race detector enabled: go test -race

	repo := NewRepository()

	// Create a channel to signal completion
	done := make(chan bool)

	// Start multiple goroutines that read and write to the repository
	for i := 0; i < 5; i++ {
		go func() {
			// Create some items
			repo.Create("Concurrent Item", "Concurrent content")
			repo.GetAll()

			// Signal completion
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify the repository has the expected number of items
	assert.Len(t, repo.items, 5)
}
