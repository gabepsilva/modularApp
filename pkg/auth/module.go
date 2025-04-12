package auth

import (
	"log"
	"os"

	"modularApp/pkg/cache"

	"github.com/gin-gonic/gin"
)

// Module represents the auth module of the application
type Module struct {
	clerkAPIKey         string
	clerkPublishableKey string
	cacheModule         *cache.Module
}

// NewModule creates a new auth module
func NewModule(cacheModule *cache.Module) *Module {
	// Load keys on initialization
	apiKey := os.Getenv("CLERK_API_KEY")
	publishableKey := os.Getenv("CLERK_PUBLISHABLE_KEY")

	module := &Module{
		clerkAPIKey:         apiKey,
		clerkPublishableKey: publishableKey,
		cacheModule:         cacheModule,
	}

	// Validate keys and log warnings if needed
	module.ValidateKeys()

	return module
}

// ValidateKeys checks if the Clerk keys are set and logs warnings if they're not
func (m *Module) ValidateKeys() {
	if m.clerkAPIKey == "" {
		log.Println("Warning: CLERK_API_KEY not set. Authentication will fail.")
	}

	if m.clerkPublishableKey == "" {
		log.Println("Warning: CLERK_PUBLISHABLE_KEY not set. UI authentication components will not work.")
	}
}

// Middleware returns the auth middleware
func (m *Module) Middleware() gin.HandlerFunc {
	var jwtCache *cache.JWTCache
	if m.cacheModule != nil {
		jwtCache = m.cacheModule.GetJWTCache()
	}
	return AuthMiddleware(m.clerkAPIKey, jwtCache)
}

// GetClerkPublishableKey returns the Clerk publishable key
func (m *Module) GetClerkPublishableKey() string {
	return m.clerkPublishableKey
}

// GetClerkAPIKey returns the Clerk API key
func (m *Module) GetClerkAPIKey() string {
	return m.clerkAPIKey
}

// GetCurrentUser retrieves the current user from the gin context
func GetCurrentUser(c *gin.Context) (User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return User{}, false
	}

	authUser, ok := user.(User)
	if !ok {
		return User{}, false
	}

	return authUser, true
}
