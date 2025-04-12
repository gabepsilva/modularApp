package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// SessionClaims simulates the Clerk SessionClaims structure
type SessionClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

// TokenCacheEntry represents a cached JWT token
type TokenCacheEntry struct {
	Claims   *SessionClaims `json:"claims"`
	UserID   string         `json:"user_id"`
	ExpireAt time.Time      `json:"expire_at"`
}

func main() {
	// Test JWT caching using Memcached
	// Using simulated JWT tokens and claims to match your production scenario

	servers := []string{
		"localhost:11211", // memcached1
		"localhost:11212", // memcached2
	}

	fmt.Println("Testing JWT token caching in Memcached...")

	// Create test tokens
	testTokens := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMSIsImV4cCI6MTY1MDAwMDAwMCwiaWF0IjoxNjAwMDAwMDAwfQ.signature",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMiIsImV4cCI6MTY1MDAwMDAwMCwiaWF0IjoxNjAwMDAwMDAwfQ.signature",
	}

	// Create test claims
	testClaims := []*SessionClaims{
		{
			Subject:   "user1",
			ExpiresAt: 1650000000,
			IssuedAt:  1600000000,
		},
		{
			Subject:   "user2",
			ExpiresAt: 1650000000,
			IssuedAt:  1600000000,
		},
	}

	// Test each server
	for i, server := range servers {
		fmt.Printf("\nTesting Memcached server #%d: %s\n", i+1, server)
		client := memcache.New(server)
		client.Timeout = time.Second * 1

		// Set TTL to 15 minutes (default in your app)
		ttl := 15 * time.Minute
		keyPrefix := "jwt:"

		// Test cache operations
		fmt.Println("Testing JWT token caching...")

		// First cache a token
		token := testTokens[0]
		claims := testClaims[0]

		entry := TokenCacheEntry{
			Claims:   claims,
			UserID:   claims.Subject,
			ExpireAt: time.Now().Add(ttl),
		}

		// Convert entry to JSON
		jsonValue, err := json.Marshal(entry)
		if err != nil {
			fmt.Printf("  ERROR: Failed to marshal cache entry: %v\n", err)
			continue
		}

		// Cache the token
		startTime := time.Now()
		err = client.Set(&memcache.Item{
			Key:        keyPrefix + token,
			Value:      jsonValue,
			Expiration: int32(ttl.Seconds()),
		})
		setTime := time.Since(startTime)

		if err != nil {
			fmt.Printf("  ERROR: Failed to cache token: %v\n", err)
			continue
		}
		fmt.Printf("  Cache set operation: %v\n", setTime)

		// Retrieve the token from cache
		startTime = time.Now()
		item, err := client.Get(keyPrefix + token)
		getTime := time.Since(startTime)

		if err != nil {
			fmt.Printf("  ERROR: Failed to retrieve token from cache: %v\n", err)
			continue
		}
		fmt.Printf("  Cache get operation: %v\n", getTime)

		// Decode the cached entry
		startTime = time.Now()
		var cachedEntry TokenCacheEntry
		err = json.Unmarshal(item.Value, &cachedEntry)
		unmarshalTime := time.Since(startTime)

		if err != nil {
			fmt.Printf("  ERROR: Failed to unmarshal cached entry: %v\n", err)
			continue
		}
		fmt.Printf("  Unmarshal operation: %v\n", unmarshalTime)
		fmt.Printf("  Cached user ID: %s\n", cachedEntry.UserID)

		// Test multiple operation performance
		fmt.Println("  Testing multiple operations (10 iterations)...")

		var totalSetTime, totalGetTime, totalUnmarshalTime time.Duration

		for j := 0; j < 10; j++ {
			// Set
			token := fmt.Sprintf("%s-%d", testTokens[j%2], j)
			claims := testClaims[j%2]

			entry := TokenCacheEntry{
				Claims:   claims,
				UserID:   claims.Subject,
				ExpireAt: time.Now().Add(ttl),
			}

			jsonValue, _ := json.Marshal(entry)

			startTime = time.Now()
			err = client.Set(&memcache.Item{
				Key:        keyPrefix + token,
				Value:      jsonValue,
				Expiration: int32(ttl.Seconds()),
			})
			opTime := time.Since(startTime)
			totalSetTime += opTime

			if err != nil {
				log.Printf("    Set error [%d]: %v", j, err)
				continue
			}

			// Get
			startTime = time.Now()
			item, err := client.Get(keyPrefix + token)
			opTime = time.Since(startTime)
			totalGetTime += opTime

			if err != nil {
				log.Printf("    Get error [%d]: %v", j, err)
				continue
			}

			// Unmarshal
			startTime = time.Now()
			var cachedEntry TokenCacheEntry
			_ = json.Unmarshal(item.Value, &cachedEntry)
			opTime = time.Since(startTime)
			totalUnmarshalTime += opTime
		}

		fmt.Printf("  Average set time: %v\n", totalSetTime/10)
		fmt.Printf("  Average get time: %v\n", totalGetTime/10)
		fmt.Printf("  Average unmarshal time: %v\n", totalUnmarshalTime/10)
		fmt.Printf("  Average total operation time: %v\n", (totalSetTime+totalGetTime+totalUnmarshalTime)/10)
	}
}
