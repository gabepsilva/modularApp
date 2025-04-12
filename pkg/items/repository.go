package items

import (
	"errors"
	"strconv"
	"sync"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

// Repository handles storage and retrieval of items
type Repository struct {
	items  map[string]Item
	mutex  sync.RWMutex
	nextID int
}

// NewRepository creates a new repository for items
func NewRepository() *Repository {
	return &Repository{
		items:  make(map[string]Item),
		nextID: 1,
	}
}

// GetAll returns all items
func (r *Repository) GetAll() []Item {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	items := make([]Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items
}

// GetByID returns an item by ID
func (r *Repository) GetByID(id string) (Item, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return Item{}, ErrItemNotFound
	}
	return item, nil
}

// Create creates a new item
func (r *Repository) Create(name, content string) Item {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := r.generateID()
	item := Item{
		ID:      id,
		Name:    name,
		Content: content,
	}
	r.items[id] = item
	return item
}

// Update updates an existing item
func (r *Repository) Update(id string, name, content string) (Item, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	item, exists := r.items[id]
	if !exists {
		return Item{}, ErrItemNotFound
	}

	if name != "" {
		item.Name = name
	}
	if content != "" {
		item.Content = content
	}
	r.items[id] = item
	return item, nil
}

// Delete removes an item
func (r *Repository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.items[id]; !exists {
		return ErrItemNotFound
	}
	delete(r.items, id)
	return nil
}

// generateID creates a unique ID for items
func (r *Repository) generateID() string {
	id := r.nextID
	r.nextID++
	return strconv.Itoa(id)
} 