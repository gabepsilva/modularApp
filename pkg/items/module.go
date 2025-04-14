package items

import (
	"log"
	"net/http"
	"time"

	"modularApp/pkg/auth"
	"modularApp/pkg/web"

	"github.com/gin-gonic/gin"
)

// Module represents the items module of the application
type Module struct {
	repo    *Repository
	web     *web.Module
	authMod *auth.Module
}

// NewModule creates a new items module
func NewModule(webModule *web.Module, authModule *auth.Module) *Module {
	return &Module{
		repo:    NewRepository(),
		web:     webModule,
		authMod: authModule,
	}
}

// RegisterRoutes registers all the routes for this module
func (m *Module) RegisterRoutes() {
	// Define base path for items API
	apiGroup := m.web.RegisterGroup("/api/v1")

	// Create a route group with authentication
	itemsGroup := apiGroup.Group("/items", m.authMod.Middleware())

	// Register item routes - all protected by auth
	itemsGroup.GET("", m.getAllItems)       // Get all items
	itemsGroup.GET("/:id", m.getItemByID)   // Get item by ID
	itemsGroup.POST("", m.createItem)       // Create a new item
	itemsGroup.PUT("/:id", m.updateItem)    // Update an existing item
	itemsGroup.DELETE("/:id", m.deleteItem) // Delete an item
}

// getAllItems returns all items
// @Summary Get all items
// @Description Retrieves a list of all items
// @Tags items
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Returns items array and user info"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/items [get]
func (m *Module) getAllItems(c *gin.Context) {
	// Get the current user (authentication is guaranteed by middleware)
	user, _ := auth.GetCurrentUser(c)

	items := m.repo.GetAll()
	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"user":  user,
	})
}

// getItemByID returns a specific item by ID
// @Summary Get item by ID
// @Description Retrieves a specific item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Item ID"
// @Success 200 {object} Item
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/items/{id} [get]
func (m *Module) getItemByID(c *gin.Context) {
	id := c.Param("id")

	item, err := m.repo.GetByID(id)

	if err != nil {
		if err == ErrItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// createItem creates a new item
// @Summary Create a new item
// @Description Creates a new item with the provided details
// @Tags items
// @Accept json
// @Produce json
// @Security Bearer
// @Param item body CreateItemRequest true "Item details"
// @Success 201 {object} map[string]interface{} "Returns created item and creator ID"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/items [post]
func (m *Module) createItem(c *gin.Context) {
	startTime := time.Now()

	// Get the current user
	getUserStart := time.Now()
	user, _ := auth.GetCurrentUser(c)
	userTime := time.Since(getUserStart)

	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repoStartTime := time.Now()
	item := m.repo.Create(req.Name, req.Content)
	repoTime := time.Since(repoStartTime)

	totalTime := time.Since(startTime)
	log.Printf("ITEMS-CREATE: total=%v, auth=%v, repo=%v", totalTime, userTime, repoTime)

	c.JSON(http.StatusCreated, gin.H{
		"item":      item,
		"createdBy": user.ID,
	})
}

// updateItem updates an existing item
// @Summary Update an existing item
// @Description Updates an item with the provided ID and details
// @Tags items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Item ID"
// @Param item body UpdateItemRequest true "Updated item details"
// @Success 200 {object} Item
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/items/{id} [put]
func (m *Module) updateItem(c *gin.Context) {
	startTime := time.Now()
	id := c.Param("id")

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repoStartTime := time.Now()
	item, err := m.repo.Update(id, req.Name, req.Content)
	repoTime := time.Since(repoStartTime)

	if err != nil {
		log.Printf("ITEMS-UPDATE: error=%v, time=%v", err, time.Since(startTime))
		if err == ErrItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalTime := time.Since(startTime)
	log.Printf("ITEMS-UPDATE: id=%s, total=%v, repo=%v", id, totalTime, repoTime)

	c.JSON(http.StatusOK, item)
}

// deleteItem deletes an item
// @Summary Delete an item
// @Description Deletes an item with the provided ID
// @Tags items
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Item ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/items/{id} [delete]
func (m *Module) deleteItem(c *gin.Context) {
	startTime := time.Now()
	id := c.Param("id")

	repoStartTime := time.Now()
	err := m.repo.Delete(id)
	repoTime := time.Since(repoStartTime)

	if err != nil {
		log.Printf("ITEMS-DELETE: error=%v, time=%v", err, time.Since(startTime))
		if err == ErrItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalTime := time.Since(startTime)
	log.Printf("ITEMS-DELETE: id=%s, total=%v, repo=%v", id, totalTime, repoTime)

	c.JSON(http.StatusNoContent, nil)
}
