package main

import (
	"fmt"
	"log"
	"os"

	"modularApp/internal/server"
)

func main() {
	// Show all environment variables for debugging
	fmt.Println("--- ENVIRONMENT VARIABLES ---")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	fmt.Println("----------------------------")

	// Check if Clerk API key is set
	clerkApiKey := os.Getenv("CLERK_API_KEY")
	if clerkApiKey == "" {
		log.Println("Warning: CLERK_API_KEY not set. Authentication will fail.")
	} else {
		keyLen := len(clerkApiKey)
		maskedKey := clerkApiKey
		if keyLen > 4 {
			maskedKey = clerkApiKey[0:4] + "..." + clerkApiKey[keyLen-4:keyLen]
		}
		log.Printf("CLERK_API_KEY found with length %d: %s", keyLen, maskedKey)
	}

	// Check if Clerk publishable key is set
	clerkPubKey := os.Getenv("CLERK_PUBLISHABLE_KEY")
	if clerkPubKey == "" {
		log.Println("Warning: CLERK_PUBLISHABLE_KEY not set. UI authentication components will not work.")
	} else {
		keyLen := len(clerkPubKey)
		maskedKey := clerkPubKey
		if keyLen > 4 {
			maskedKey = clerkPubKey[0:4] + "..." + clerkPubKey[keyLen-4:keyLen]
		}
		log.Printf("CLERK_PUBLISHABLE_KEY found with length %d: %s", keyLen, maskedKey)
	}

	// Create server with all modules
	srv := server.New()

	// Start the server
	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
