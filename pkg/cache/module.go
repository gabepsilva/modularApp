package cache

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Module represents the cache module of the application
type Module struct {
	jwtCache *JWTCache
}

// NewModule creates a new cache module
func NewModule() *Module {
	// Get memcached servers from environment
	memcachedServers := getMemcachedServers()

	// Get JWT cache TTL from environment or use default (15 minutes)
	jwtCacheTTL := getJWTCacheTTL()

	return &Module{
		jwtCache: NewJWTCache(memcachedServers, "jwt:", jwtCacheTTL),
	}
}

// GetJWTCache returns the JWT cache
func (m *Module) GetJWTCache() *JWTCache {
	return m.jwtCache
}

// getMemcachedServers returns a list of memcached servers from environment
// or default servers if not set
func getMemcachedServers() []string {
	// Default servers based on docker-compose configuration
	defaultServers := []string{
		"localhost:11211", // memcached1
		"localhost:11212", // memcached2
	}

	// Get servers from environment
	serversEnv := os.Getenv("MEMCACHED_SERVERS")
	if serversEnv == "" {
		return defaultServers
	}

	return strings.Split(serversEnv, ",")
}

// getJWTCacheTTL returns the JWT cache TTL from environment or default
func getJWTCacheTTL() time.Duration {
	// Default TTL: 15 minutes
	defaultTTL := 15 * time.Minute

	// Get TTL from environment
	ttlEnv := os.Getenv("JWT_CACHE_TTL_SECONDS")
	if ttlEnv == "" {
		return defaultTTL
	}

	// Parse TTL
	ttl, err := strconv.Atoi(ttlEnv)
	if err != nil {
		return defaultTTL
	}

	return time.Duration(ttl) * time.Second
}
