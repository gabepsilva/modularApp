package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

// User contains basic user information from Clerk
type User struct {
	ID    string
	Email string
	Name  string
}

// AuthMiddleware creates a gin middleware for Clerk authentication
func AuthMiddleware() gin.HandlerFunc {
	// Get Clerk API key from environment variable
	clerkAPIKey := os.Getenv("CLERK_API_KEY")
	if clerkAPIKey == "" {
		// If the API key is not set, return a middleware that always fails
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Clerk API key not configured"})
			c.Abort()
			return
		}
	}

	// Create Clerk client once when middleware is initialized
	client, err := clerk.NewClient(clerkAPIKey)
	if err != nil {
		// If client creation fails, return a middleware that always fails
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Clerk client"})
			c.Abort()
			return
		}
	}

	// Return the actual middleware function with the pre-initialized client
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the "Bearer " prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer token"})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Verify the session token using the pre-initialized client
		claims, err := client.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		// Get user details from Clerk using the pre-initialized client
		userId := claims.Subject
		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token does not contain a user ID"})
			c.Abort()
			return
		}

		user, err := client.Users().Read(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get user details: %v", err)})
			c.Abort()
			return
		}

		// For safety, handle nil fields
		emailAddress := ""
		if len(user.EmailAddresses) > 0 && user.EmailAddresses[0].EmailAddress != "" {
			emailAddress = user.EmailAddresses[0].EmailAddress
		}

		firstName := ""
		if user.FirstName != nil {
			firstName = *user.FirstName
		}

		lastName := ""
		if user.LastName != nil {
			lastName = *user.LastName
		}

		// Store user in context for use in handlers
		c.Set("user", User{
			ID:    user.ID,
			Email: emailAddress,
			Name:  firstName + " " + lastName,
		})

		c.Next()
	}
}
