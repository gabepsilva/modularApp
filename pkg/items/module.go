package items

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"modularApp/pkg/web"
)

// Module represents the items module of the application
type Module struct {
	repo *Repository
	web  *web.Module
}

// NewModule creates a new items module
func NewModule(webModule *web.Module) *Module {
	return &Module{
		repo: NewRepository(),
		web:  webModule,
	}
}

// RegisterRoutes registers all the routes for this module
func (m *Module) RegisterRoutes() {
	// Define base path for items API
	apiGroup := m.web.RegisterGroup("/api/v1")
	itemsGroup := apiGroup.Group("/items")
	
	// Register item routes
	itemsGroup.GET("", m.getAllItems)
	itemsGroup.GET("/:id", m.getItemByID)
	itemsGroup.POST("", m.createItem)
	itemsGroup.PUT("/:id", m.updateItem)
	itemsGroup.DELETE("/:id", m.deleteItem)
}

// getAllItems returns all items
func (m *Module) getAllItems(c *gin.Context) {
	items := m.repo.GetAll()
	c.JSON(http.StatusOK, items)
}

// getItemByID returns a specific item by ID
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
func (m *Module) createItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	item := m.repo.Create(req.Name, req.Content)
	c.JSON(http.StatusCreated, item)
}

// updateItem updates an existing item
func (m *Module) updateItem(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	item, err := m.repo.Update(id, req.Name, req.Content)
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

// deleteItem deletes an item
func (m *Module) deleteItem(c *gin.Context) {
	id := c.Param("id")
	
	err := m.repo.Delete(id)
	if err != nil {
		if err == ErrItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusNoContent, nil)
} 