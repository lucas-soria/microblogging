package main

import (
	"log"

	"github.com/lucas-soria/microblogging/cmd/users/handlers"

	"github.com/lucas-soria/microblogging/internal/users"
)

func main() {
	// Initialize repository
	log.Println("Initializing users repository")
	userRepo := users.NewInMemoryUserRepository()

	// Initialize service with repository
	log.Println("Initializing users service")
	userService := users.NewUserService(userRepo)

	// Initialize handlers with service
	log.Println("Initializing users handlers")
	userHandler := handlers.NewUserHandler(userService)

	// Create application
	log.Println("Creating users application")
	application := NewApplication(userHandler)

	server := newServer()

	addRoutes(server, application)

	// Start server
	server.Start()
}
