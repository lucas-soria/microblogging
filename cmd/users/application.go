package main

import (
	"github.com/lucas-soria/microblogging/cmd/users/handlers"
)

// Application holds the dependencies for the HTTP server
type Application struct {
	userHandler *handlers.UserHandler
}

// NewApplication creates a new HTTP server and sets up routing
func NewApplication(userHandler *handlers.UserHandler) *Application {
	application := &Application{
		userHandler: userHandler,
	}

	return application
}
