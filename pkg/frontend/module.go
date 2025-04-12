package frontend

import (
	"net/http"

	"modularApp/pkg/auth"
	"modularApp/pkg/web"

	"github.com/gin-gonic/gin"
)

// Module represents the frontend module
type Module struct {
	web  *web.Module
	auth *auth.Module
}

// NewModule creates a new frontend module
func NewModule(webModule *web.Module, authModule *auth.Module) *Module {
	// Set HTML renderer on the router
	webModule.Router().LoadHTMLGlob("web/templates/*.html")

	return &Module{
		web:  webModule,
		auth: authModule,
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
		"clerkPublishableKey": m.auth.GetClerkPublishableKey(),
		"template":            "home.html",
	})
}

// loginHandler handles the login page request
func (m *Module) loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "layout", gin.H{
		"title":               "Login - Modular App",
		"clerkPublishableKey": m.auth.GetClerkPublishableKey(),
		"template":            "login.html",
	})
}

// dashboardHandler handles the dashboard page request
func (m *Module) dashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "layout", gin.H{
		"title":               "Dashboard - Modular App",
		"clerkPublishableKey": m.auth.GetClerkPublishableKey(),
		"template":            "dashboard.html",
	})
}
