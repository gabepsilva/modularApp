package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"modularApp/pkg/cache"

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
func AuthMiddleware(clerkAPIKey string, jwtCache *cache.JWTCache) gin.HandlerFunc {
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

		// Try to get the token from cache first
		var claims *clerk.SessionClaims
		var cachedUser *cache.User
		var err error
		cacheHit := false
		userCacheHit := false

		if jwtCache != nil {
			// Check if token is in cache
			cachedClaims, user, found := jwtCache.GetToken(token)

			if found {
				claims = cachedClaims
				cachedUser = user
				cacheHit = true
				userCacheHit = (user != nil)
			}
		}

		// If not in cache, verify the token
		if claims == nil {
			claims, err = client.VerifyToken(token)

			if err != nil {
				log.Printf("TOKEN VERIFICATION FAILED: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
				c.Abort()
				return
			}
		}

		// Get user ID from claims
		userId := claims.Subject
		if userId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token does not contain a user ID"})
			c.Abort()
			return
		}

		// User object to store in context
		var user User

		// Get user details from cache or API
		if userCacheHit && cachedUser != nil {
			// Use cached user data
			user = User{
				ID:    cachedUser.ID,
				Email: cachedUser.Email,
				Name:  cachedUser.Name,
			}
		} else {
			// Get user details from Clerk API
			clerkUser, err := client.Users().Read(userId)

			if err != nil {
				log.Printf("USER FETCH FAILED: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get user details: %v", err)})
				c.Abort()
				return
			}

			// For safety, handle nil fields
			emailAddress := ""
			if len(clerkUser.EmailAddresses) > 0 && clerkUser.EmailAddresses[0].EmailAddress != "" {
				emailAddress = clerkUser.EmailAddresses[0].EmailAddress
			}

			firstName := ""
			if clerkUser.FirstName != nil {
				firstName = *clerkUser.FirstName
			}

			lastName := ""
			if clerkUser.LastName != nil {
				lastName = *clerkUser.LastName
			}

			// Create user object
			user = User{
				ID:    clerkUser.ID,
				Email: emailAddress,
				Name:  firstName + " " + lastName,
			}

			// Cache user data for future requests
			if jwtCache != nil && !cacheHit {
				cacheUser := &cache.User{
					ID:    user.ID,
					Email: user.Email,
					Name:  user.Name,
				}

				err := jwtCache.StoreToken(token, claims, cacheUser)
				if err != nil {
					log.Printf("CACHE STORE FAILED: %v", err)
				}
			}
		}

		// Store user in context for use in handlers
		c.Set("user", user)
		c.Next()
	}
}
