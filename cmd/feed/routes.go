package main

import (
	"github.com/gin-gonic/gin"
)

// setupRoutes configures all the routes for the application
func addRoutes(server *Server, application *Application) {
	healthCheck(server.router)
	groupV1 := newGroup(server.router, "/v1")
	feedRoutes(groupV1, application)
}

func newGroup(router *gin.Engine, path string) *gin.RouterGroup {
	return router.Group(path)
}

func healthCheck(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func feedRoutes(group *gin.RouterGroup, application *Application) {
	group.GET("/feed/timeline", application.feedHandler.GetUserTimeline)
}
