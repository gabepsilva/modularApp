package items

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"modularApp/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup test environment for testing the handlers directly
func setupItemsTest() (*gin.Engine, *Repository) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router and repository
	router := gin.New()
	repo := NewRepository()

	// Create a basic API group for testing
	apiGroup := router.Group("/api/v1")

	// Setup item routes with a simple auth middleware that just adds a test user
	itemsGroup := apiGroup.Group("/items", func(c *gin.Context) {
		// Mock auth middleware that just sets a user
		c.Set("user", auth.User{
			ID:    "test-user-id",
			Email: "test@example.com",
			Name:  "Test User",
		})
		c.Next()
	})

	// Create a simple module with just the repository for handler testing
	module := &Module{repo: repo}

	// Register the routes directly
	itemsGroup.GET("", module.getAllItems)
	itemsGroup.GET("/:id", module.getItemByID)
	itemsGroup.POST("", module.createItem)
	itemsGroup.PUT("/:id", module.updateItem)
	itemsGroup.DELETE("/:id", module.deleteItem)

	return router, repo
}

// Helper function to convert an object to JSON
func toJSON(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func TestGetAllItems(t *testing.T) {
	// Setup
	router, repo := setupItemsTest()

	// Add some test data
	repo.Create("Test Item 1", "Content 1")
	repo.Create("Test Item 2", "Content 2")

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/items", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Check that items array exists and has 2 items
	items, exists := response["items"].([]interface{})
	assert.True(t, exists)
	assert.Len(t, items, 2)

	// Verify user info is present
	user, exists := response["user"].(map[string]interface{})
	assert.True(t, exists)
	assert.Equal(t, "test-user-id", user["ID"])
}

func TestGetItemByID(t *testing.T) {
	// Setup
	router, repo := setupItemsTest()

	// Add a test item
	item := repo.Create("Test Item", "Test Content")

	// Create a test request for existing item
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/items/"+item.ID, nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response body
	var responseItem Item
	json.Unmarshal(w.Body.Bytes(), &responseItem)
	assert.Equal(t, item.ID, responseItem.ID)
	assert.Equal(t, item.Name, responseItem.Name)
	assert.Equal(t, item.Content, responseItem.Content)

	// Test non-existent item
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/items/999", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response for non-existent item
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateItem(t *testing.T) {
	// Setup
	router, _ := setupItemsTest()

	// Create a test request body
	createReq := CreateItemRequest{
		Name:    "New Item",
		Content: "New Content",
	}
	body, _ := json.Marshal(createReq)

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Check item details
	item, exists := response["item"].(map[string]interface{})
	assert.True(t, exists)
	assert.Equal(t, "New Item", item["name"])
	assert.Equal(t, "New Content", item["content"])

	// Check createdBy field
	assert.Equal(t, "test-user-id", response["createdBy"])

	// Test invalid request body
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/items", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response for invalid request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateItem(t *testing.T) {
	// Setup
	router, repo := setupItemsTest()

	// Add a test item
	item := repo.Create("Original Name", "Original Content")

	// Create update request
	updateReq := UpdateItemRequest{
		Name:    "Updated Name",
		Content: "Updated Content",
	}
	body, _ := json.Marshal(updateReq)

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/items/"+item.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response body
	var responseItem Item
	json.Unmarshal(w.Body.Bytes(), &responseItem)
	assert.Equal(t, item.ID, responseItem.ID)
	assert.Equal(t, "Updated Name", responseItem.Name)
	assert.Equal(t, "Updated Content", responseItem.Content)

	// Test updating non-existent item
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/items/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response for non-existent item
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteItem(t *testing.T) {
	// Setup
	router, repo := setupItemsTest()

	// Add a test item
	item := repo.Create("Test Item", "Test Content")

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/items/"+item.ID, nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify item was deleted
	items := repo.GetAll()
	assert.Empty(t, items)

	// Test deleting non-existent item
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/items/999", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response for non-existent item
	assert.Equal(t, http.StatusNotFound, w.Code)
}
