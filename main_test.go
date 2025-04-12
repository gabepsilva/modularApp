package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"modularApp/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareIntegration(t *testing.T) {
	// Set test environment variables
	os.Setenv("CLERK_API_KEY", "test-api-key")
	os.Setenv("CLERK_PUBLISHABLE_KEY", "test-publishable-key")

	// Create a new gin engine
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create auth module
	authModule := auth.NewModule()

	// Setup a protected route
	r.GET("/protected", authModule.Middleware(), func(c *gin.Context) {
		user, exists := auth.GetCurrentUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "This is a protected route",
			"user_id": user.ID,
		})
	})

	// Test with missing Authorization header
	t.Run("Missing Authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	// Test with invalid token
	t.Run("Invalid token format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header must be Bearer token")
	})

	// Note: We can't test a successful authentication in an integration test
	// without mocking the Clerk HTTP calls, which is more complex.
	// For a complete test, you would need to use a mocking library for HTTP requests
	// or create a test-specific middleware that accepts certain test tokens.
}

// TestMain is the entry point for all tests in the main package
func TestMain(m *testing.M) {
	// Setup before tests
	gin.SetMode(gin.TestMode)

	// Run tests
	exitCode := m.Run()

	// Cleanup after tests

	// Exit with the same code as the tests
	os.Exit(exitCode)
}
