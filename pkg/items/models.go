package items

// Item represents a basic resource item
type Item struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// CreateItemRequest is used for creating new items
type CreateItemRequest struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UpdateItemRequest is used for updating existing items
type UpdateItemRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
} 