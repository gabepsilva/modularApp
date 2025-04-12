package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

// User contains basic user information from Clerk
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TokenCacheEntry represents a cached JWT token
type TokenCacheEntry struct {
	Claims   *clerk.SessionClaims `json:"claims"`
	UserID   string               `json:"user_id"`
	User     *User                `json:"user,omitempty"`
	ExpireAt time.Time            `json:"expire_at"`
}

// JWTCache is a specialized cache for JWT tokens
type JWTCache struct {
	cache     *MemcachedClient
	keyPrefix string
	ttl       time.Duration
}

// NewJWTCache creates a new JWT token cache
func NewJWTCache(memcachedServers []string, keyPrefix string, ttl time.Duration) *JWTCache {
	return &JWTCache{
		cache:     NewMemcachedClient(memcachedServers),
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

// hashKey creates a safe key for Memcached by hashing the token
func (j *JWTCache) hashKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return j.keyPrefix + hex.EncodeToString(hash[:])
}

// StoreToken caches a JWT token and its validation result
func (j *JWTCache) StoreToken(token string, claims *clerk.SessionClaims, user *User) error {
	// Create a cache entry
	entry := TokenCacheEntry{
		Claims:   claims,
		UserID:   claims.Subject,
		User:     user,
		ExpireAt: time.Now().Add(j.ttl),
	}

	// Store in cache with hashed key
	return j.cache.Set(j.hashKey(token), entry, j.ttl)
}

// GetToken retrieves a cached JWT token validation result
func (j *JWTCache) GetToken(token string) (*clerk.SessionClaims, *User, bool) {
	var entry TokenCacheEntry
	err := j.cache.Get(j.hashKey(token), &entry)
	if err != nil {
		return nil, nil, false
	}

	// Check if the entry is still valid
	if time.Now().After(entry.ExpireAt) {
		// Token is expired, remove it from cache
		j.cache.Delete(j.hashKey(token))
		return nil, nil, false
	}

	return entry.Claims, entry.User, true
}

// InvalidateToken removes a token from the cache
func (j *JWTCache) InvalidateToken(token string) error {
	return j.cache.Delete(j.hashKey(token))
}
