package main

import (
	"github.com/lucas-soria/microblogging/cmd/tweets/handlers"
)

// Application holds the dependencies for the HTTP server
type Application struct {
	tweetHandler *handlers.TweetHandler
}

// NewApplication creates a new HTTP server and sets up routing
func NewApplication(tweetHandler *handlers.TweetHandler) *Application {
	application := &Application{
		tweetHandler: tweetHandler,
	}

	return application
}
