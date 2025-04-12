package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	// Test Memcached key hashing for JWT tokens

	// Sample JWT token that would be too long as a direct key
	longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkwMjIsImF1ZCI6InNvbWUtYXVkaWVuY2UiLCJpc3MiOiJzb21lLWlzc3VlciIsImp0aSI6InNvbWUtdW5pcXVlLWlkIiwibmJmIjoxNTE2MjM5MDIyLCJ0eXAiOiJzb21lLXR5cGUifQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Original key approach (will fail)
	originalKey := "jwt:" + longToken
	fmt.Printf("Original key length: %d characters\n", len(originalKey))

	// Hashed key approach
	hash := sha256.Sum256([]byte(longToken))
	hashedKey := "jwt:" + hex.EncodeToString(hash[:])
	fmt.Printf("Hashed key length: %d characters\n", len(hashedKey))
	fmt.Printf("Hashed key: %s\n", hashedKey)

	// Test with Memcached
	client := memcache.New("localhost:11211")
	client.Timeout = time.Second * 1

	// Try to set with original key (should fail)
	err := client.Set(&memcache.Item{
		Key:        originalKey,
		Value:      []byte("test value"),
		Expiration: 60,
	})
	if err != nil {
		fmt.Printf("Original key error: %v\n", err)
	} else {
		fmt.Println("Original key worked (unexpected)")
	}

	// Try to set with hashed key (should work)
	err = client.Set(&memcache.Item{
		Key:        hashedKey,
		Value:      []byte("test value"),
		Expiration: 60,
	})
	if err != nil {
		fmt.Printf("Hashed key error: %v\n", err)
	} else {
		fmt.Println("Hashed key worked as expected")

		// Try to get with hashed key
		item, err := client.Get(hashedKey)
		if err != nil {
			fmt.Printf("Hashed key get error: %v\n", err)
		} else {
			fmt.Printf("Successfully retrieved with hashed key: %s\n", string(item.Value))
		}
	}
}
