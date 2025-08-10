package main

import (
	"log"

	"github.com/lucas-soria/microblogging/cmd/tweets/handlers"

	"github.com/lucas-soria/microblogging/internal/tweets"
)

func main() {
	// Initialize repository
	log.Println("Initializing tweets repository")
	tweetRepo := tweets.NewInMemoryTweetRepository()

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
