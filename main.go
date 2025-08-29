package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alnav3/splitweb/routes"
	"github.com/joho/godotenv"
)


func main() {

	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Set up all routes
	routes.SetupRoutes()



	// Start the server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// Helper functions
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
