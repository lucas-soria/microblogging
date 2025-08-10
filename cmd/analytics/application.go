package main

import (
	"github.com/lucas-soria/microblogging/cmd/analytics/handlers"
)

// Application holds the dependencies for the HTTP server
type Application struct {
	analyticsHandler *handlers.AnalyticsHandler
}

// NewApplication creates a new HTTP server and sets up routing
func NewApplication(analyticsHandler *handlers.AnalyticsHandler) *Application {
	application := &Application{
		analyticsHandler: analyticsHandler,
	}

	return application
}
