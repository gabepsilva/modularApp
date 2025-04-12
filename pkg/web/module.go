package web

import (
	"github.com/gin-gonic/gin"
)

// Module represents the web module of the application
type Module struct {
	router *gin.Engine
}

// NewModule creates a new web module
func NewModule() *Module {
	router := gin.Default()
	return &Module{
		router: router,
	}
}

// Router returns the underlying router instance
func (m *Module) Router() *gin.Engine {
	return m.router
}

// RegisterRoutes registers the base routes for this module
func (m *Module) RegisterRoutes() {
	// Root group for all web module routes
	group := m.router.Group("/api")

	// Define routes
	group.GET("/health", m.healthCheck)

	// You can add more route groups
	v1 := group.Group("/v1")
	{
		v1.GET("/ping", m.ping)
	}
}

// RegisterHandler allows other modules to register their handlers
func (m *Module) RegisterHandler(path string, method string, handlers ...gin.HandlerFunc) {
	m.router.Handle(method, path, handlers...)
}

// RegisterGroup creates a new route group and allows other modules to register their group handlers
func (m *Module) RegisterGroup(path string) *gin.RouterGroup {
	return m.router.Group(path)
}

// healthCheck handler for health check endpoint
func (m *Module) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
	})
}

// ping handler for ping endpoint
func (m *Module) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
