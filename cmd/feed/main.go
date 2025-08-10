package main

import (
	"log"

	"github.com/lucas-soria/microblogging/cmd/feed/handlers"

	"github.com/lucas-soria/microblogging/internal/feed"
)

func main() {
	// Initialize repository
	log.Println("Initializing feed repository")
	feedRepo := feed.NewInMemoryFeedRepository()

	// Initialize service with repository
	log.Println("Initializing feed service")
	feedService := feed.NewService(feedRepo)

	// Initialize handlers with service
	log.Println("Initializing feed handlers")
	feedHandler := handlers.NewFeedHandler(feedService)

	// Create application
	log.Println("Creating feed application")
	application := NewApplication(feedHandler)

	server := newServer()

	addRoutes(server, application)

	// Start server
	server.Start()
}
