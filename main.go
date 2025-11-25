package main

import (
	"log"
	"net/http"
	"os"

	databaselogic "github.com/alnav3/splitweb/db/database_logic"
	"github.com/alnav3/splitweb/routes"
	"github.com/joho/godotenv"
)


func main() {

	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Setup database
	repo, err := databaselogic.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer repo.Close()

	log.Println("Database setup completed successfully")

	// Set up all routes
	routes.SetupRoutes(repo)

	// Start the server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	err = http.ListenAndServe(":"+port, nil)
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
