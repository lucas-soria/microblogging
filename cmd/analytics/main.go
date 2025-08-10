package main

import (
	"log"

	"github.com/lucas-soria/microblogging/cmd/analytics/handlers"

	"github.com/lucas-soria/microblogging/internal/analytics"
)

func main() {
	// Initialize repository
	log.Println("Initializing feed repository")
	analyticsRepo := analytics.NewInMemoryRepository()

	// Initialize service with repository
	log.Println("Initializing feed service")
	analyticsService := analytics.NewAnalyticsService(analyticsRepo)

	// Initialize handlers with service
	log.Println("Initializing feed handlers")
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Create application
	log.Println("Creating feed application")
	application := NewApplication(analyticsHandler)

	server := newServer()

	addRoutes(server, application)

	// Start server
	server.Start()
}
