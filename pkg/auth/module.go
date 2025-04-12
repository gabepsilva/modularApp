package auth

import (
	"github.com/gin-gonic/gin"
)

// Module represents the auth module of the application
type Module struct{}

// NewModule creates a new auth module
func NewModule() *Module {
	return &Module{}
}

// Middleware returns the auth middleware
func (m *Module) Middleware() gin.HandlerFunc {
	return AuthMiddleware()
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
