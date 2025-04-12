package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	// Test memcached servers
	servers := []string{
		"localhost:11211", // memcached1
		"localhost:11212", // memcached2
	}

	fmt.Println("Testing Memcached connectivity...")

	// Test each server individually
	for i, server := range servers {
		fmt.Printf("Testing Memcached server #%d: %s\n", i+1, server)
		client := memcache.New(server)
		client.Timeout = time.Second * 1

		// Test connectivity
		startTime := time.Now()
		err := client.Set(&memcache.Item{
			Key:        fmt.Sprintf("test-key-%d", i),
			Value:      []byte("test-value"),
			Expiration: 60,
		})
		setTime := time.Since(startTime)

		if err != nil {
			fmt.Printf("  ERROR: Could not connect to Memcached server: %v\n", err)
			continue
		}
		fmt.Printf("  Set operation: %v\n", setTime)

		// Test get operation
		startTime = time.Now()
		item, err := client.Get(fmt.Sprintf("test-key-%d", i))
		getTime := time.Since(startTime)

		if err != nil {
			fmt.Printf("  ERROR: Could not get key from Memcached server: %v\n", err)
			continue
		}
		fmt.Printf("  Get operation: %v\n", getTime)
		fmt.Printf("  Value: %s\n", string(item.Value))

		// Test latency with multiple operations
		fmt.Println("  Testing latency with 10 operations...")

		var totalSetTime, totalGetTime time.Duration

		for j := 0; j < 10; j++ {
			key := fmt.Sprintf("test-key-%d-%d", i, j)

			// Set
			startTime = time.Now()
			err = client.Set(&memcache.Item{
				Key:        key,
				Value:      []byte(fmt.Sprintf("test-value-%d", j)),
				Expiration: 60,
			})
			opTime := time.Since(startTime)
			totalSetTime += opTime

			if err != nil {
				log.Printf("    Set error [%d]: %v", j, err)
				continue
			}

			// Get
			startTime = time.Now()
			_, err = client.Get(key)
			opTime = time.Since(startTime)
			totalGetTime += opTime

			if err != nil {
				log.Printf("    Get error [%d]: %v", j, err)
			}
		}

		fmt.Printf("  Average set time: %v\n", totalSetTime/10)
		fmt.Printf("  Average get time: %v\n", totalGetTime/10)
	}
}
