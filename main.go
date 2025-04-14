package main

import (
	"log"

	"modularApp/internal/server"

	_ "modularApp/docs" // Import generated docs
)

// @title Modular App API
// @version 1.0
// @description A modular Go API with authentication
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345".
func main() {
	// Create server with all modules (Clerk keys are now read in auth module)
	srv := server.New()

	// Start the server
	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
