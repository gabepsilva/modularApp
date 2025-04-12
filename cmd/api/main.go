package main

import (
	"log"
	"modularApp/internal/server"
)

func main() {
	// Create server with all modules
	srv := server.New()
	
	// Start the server
	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 