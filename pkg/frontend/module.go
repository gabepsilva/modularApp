package frontend

import (
	"fmt"
	"net/http"
	"os"

	"modularApp/pkg/web"

	"github.com/gin-gonic/gin"
)

// Module represents the frontend module
type Module struct {
	web *web.Module
}

// NewModule creates a new frontend module
func NewModule(webModule *web.Module) *Module {
	// Set HTML renderer on the router
	webModule.Router().LoadHTMLGlob("web/templates/*.html")

	// Debug log to confirm templates are loaded
	fmt.Println("HTML templates loaded for frontend module")

	return &Module{
		web: webModule,
	}
}

// RegisterRoutes registers the frontend routes
func (m *Module) RegisterRoutes() {
	// Create a route group for frontend
	router := m.web.Router()

	// Routes for HTML pages
	router.GET("/", m.homeHandler)
	router.GET("/login", m.loginHandler)
	router.GET("/dashboard", m.dashboardHandler)
}

// homeHandler handles the home page request
func (m *Module) homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "layout", gin.H{
		"title":               "Modular App",
		"clerkPublishableKey": getClerkPublishableKey(),
		"template":            "home.html",
	})
}

// loginHandler handles the login page request
func (m *Module) loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "layout", gin.H{
		"title":               "Login - Modular App",
		"clerkPublishableKey": getClerkPublishableKey(),
		"template":            "login.html",
	})
}

// dashboardHandler handles the dashboard page request
func (m *Module) dashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "layout", gin.H{
		"title":               "Dashboard - Modular App",
		"clerkPublishableKey": getClerkPublishableKey(),
		"template":            "dashboard.html",
	})
}

// getClerkPublishableKey returns the Clerk publishable key from environment variables
func getClerkPublishableKey() string {
	key := os.Getenv("CLERK_PUBLISHABLE_KEY")
	fmt.Printf("DEBUG - CLERK_PUBLISHABLE_KEY value: '%s' (length: %d)\n", key, len(key))
	return key
}
