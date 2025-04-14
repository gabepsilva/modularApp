package web

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Module represents the web module of the application
type Module struct {
	router *gin.Engine
}

// NewModule creates a new web module
func NewModule() *Module {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./web/static")

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

	// Swagger documentation route
	m.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
// @Summary Server health check
// @Description Check if the server is running properly
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "Status healthy"
// @Router /api/health [get]
func (m *Module) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
	})
}

// ping handler for ping endpoint
// @Summary Ping test
// @Description Simple ping response for testing
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "Pong response"
// @Router /api/v1/ping [get]
func (m *Module) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
