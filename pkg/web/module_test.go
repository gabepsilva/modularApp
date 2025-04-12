package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set Gin to test mode for all tests
	gin.SetMode(gin.TestMode)
}

// TestNewModule tests the creation of a new web module
func TestNewModule(t *testing.T) {
	module := NewModule()

	// Verify the module was created successfully
	assert.NotNil(t, module)
	assert.NotNil(t, module.router)
}

// TestRouter tests that Router() returns the underlying router
func TestRouter(t *testing.T) {
	module := NewModule()

	// Get the router
	router := module.Router()

	// Verify the router is the same as the one in the module
	assert.NotNil(t, router)
	assert.Equal(t, module.router, router)
}

// TestRegisterRoutes tests that RegisterRoutes sets up the expected endpoints
func TestRegisterRoutes(t *testing.T) {
	module := NewModule()

	// Register routes
	module.RegisterRoutes()

	// Test the health endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	module.router.ServeHTTP(w, req)

	// Verify response from health endpoint
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "healthy", response["status"])

	// Test the ping endpoint
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/ping", nil)
	module.router.ServeHTTP(w, req)

	// Verify response from ping endpoint
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "pong", response["message"])
}

// TestRegisterHandler tests that RegisterHandler adds a handler to the router
func TestRegisterHandler(t *testing.T) {
	module := NewModule()

	// Register a test handler
	module.RegisterHandler("/test", "GET", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"result": "test_success",
		})
	})

	// Test the registered handler
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	module.router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "test_success", response["result"])
}

// TestRegisterGroup tests that RegisterGroup creates a route group
func TestRegisterGroup(t *testing.T) {
	module := NewModule()

	// Register a test group
	group := module.RegisterGroup("/test-group")
	assert.NotNil(t, group)

	// Add a handler to the group
	group.GET("/endpoint", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"group":    "test-group",
			"endpoint": "endpoint",
		})
	})

	// Test the group handler
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-group/endpoint", nil)
	module.router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "test-group", response["group"])
	assert.Equal(t, "endpoint", response["endpoint"])
}

// TestHealthCheckHandler tests the health check handler directly
func TestHealthCheckHandler(t *testing.T) {
	module := NewModule()

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	module.healthCheck(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "healthy", response["status"])
}

// TestPingHandler tests the ping handler directly
func TestPingHandler(t *testing.T) {
	module := NewModule()

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	module.ping(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "pong", response["message"])
}

// TestStaticFilesRoute tests that the static files route is registered
func TestStaticFilesRoute(t *testing.T) {
	// This test is limited since we can't easily test static file serving in unit tests
	// We're just checking that the route is registered and doesn't cause unexpected errors

	module := NewModule()

	// Test access to a nonexistent static file
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/static/nonexistent.css", nil)
	module.router.ServeHTTP(w, req)

	// The response code could be 404 (file not found) or potentially
	// another code if the directory itself doesn't exist in the test environment
	// We'll just verify that the server responded and didn't crash
	// Status should be either 404 (file not found) or 301 (directory redirect)
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusMovedPermanently,
		"Expected 404 Not Found or 301 Moved Permanently, got %d", w.Code)
}
