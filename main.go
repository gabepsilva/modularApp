package main

import (
	"log"
	"os"

	"modularApp/internal/server"
)

func main() {
	// Check if Clerk API key is set
	clerkApiKey := os.Getenv("CLERK_API_KEY")
	if clerkApiKey == "" {
		log.Println("Warning: CLERK_API_KEY not set. Authentication will fail.")
	}

	// Check if Clerk publishable key is set
	clerkPubKey := os.Getenv("CLERK_PUBLISHABLE_KEY")
	if clerkPubKey == "" {
		log.Println("Warning: CLERK_PUBLISHABLE_KEY not set. UI authentication components will not work.")
	}

	// Create server with all modules
	srv := server.New()

	// Start the server
	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
