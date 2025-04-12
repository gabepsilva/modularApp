package cache

import (
	"encoding/json"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedClient represents a load-balanced Memcached client
type MemcachedClient struct {
	clients []*memcache.Client
}

// NewMemcachedClient creates a new load-balanced Memcached client
func NewMemcachedClient(servers []string) *MemcachedClient {
	clients := make([]*memcache.Client, len(servers))
	for i, server := range servers {
		clients[i] = memcache.New(server)
	}
	return &MemcachedClient{
		clients: clients,
	}
}

// Set stores an item in cache with the given key, value and expiration time
func (m *MemcachedClient) Set(key string, value interface{}, expiration time.Duration) error {
	// Convert value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Create Memcached item
	item := &memcache.Item{
		Key:        key,
		Value:      jsonValue,
		Expiration: int32(expiration.Seconds()),
	}

	// Try each client until successful or all fail
	var lastError error
	for _, client := range m.clients {
		err := client.Set(item)
		if err == nil {
			return nil
		}
		lastError = err
	}

	return lastError
}

// Get retrieves an item from cache with the given key
func (m *MemcachedClient) Get(key string, value interface{}) error {
	// Try each client until successful or all fail
	var lastError error
	for _, client := range m.clients {
		item, err := client.Get(key)
		if err == nil {
			// Unmarshal JSON value
			return json.Unmarshal(item.Value, value)
		}
		lastError = err
	}

	return lastError
}

// Delete removes an item from cache with the given key
func (m *MemcachedClient) Delete(key string) error {
	// Try each client until successful or all fail
	var lastError error
	for _, client := range m.clients {
		err := client.Delete(key)
		if err == nil {
			return nil
		}
		lastError = err
	}

	return lastError
}
