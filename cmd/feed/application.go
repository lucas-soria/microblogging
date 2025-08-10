package main

import (
	"github.com/lucas-soria/microblogging/cmd/feed/handlers"
)

// Application holds the dependencies for the HTTP server
type Application struct {
	feedHandler *handlers.FeedHandler
}

// NewApplication creates a new HTTP server and sets up routing
func NewApplication(feedHandler *handlers.FeedHandler) *Application {
	application := &Application{
		feedHandler: feedHandler,
	}

	return application
}
