package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lucas-soria/microblogging/cmd/tweets/handlers"

	"github.com/lucas-soria/microblogging/internal/tweets"

	"github.com/lucas-soria/microblogging/pkg/database"
)

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "postgres-primary")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "not-found")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Println("Initializing tweets database connection")
	db, err := database.NewPostgresClient(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize tweets database: %v", err)
	}

	// Initialize repository
	log.Println("Initializing tweets repository")
	tweetRepo := tweets.NewPostgresTweetRepository(db)

	// Initialize service with repository
	log.Println("Initializing tweets service")
	tweetService := tweets.NewService(tweetRepo)

	// Initialize handlers with service
	log.Println("Initializing tweets handlers")
	tweetHandler := handlers.NewTweetHandler(tweetService)

	// Create application
	log.Println("Creating tweets application")
	application := NewApplication(tweetHandler)

	server := newServer()

	addRoutes(server, application)

	// Start server
	server.Start()
}
