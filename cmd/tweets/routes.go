package main

import (
	"github.com/gin-gonic/gin"
)

// setupRoutes configures all the routes for the application
func addRoutes(server *Server, application *Application) {
	healthCheck(server.router)
	groupV1 := newGroup(server.router, "/v1")
	tweetsRoutes(groupV1, application)
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

func tweetsRoutes(group *gin.RouterGroup, application *Application) {
	group.POST("/tweets", application.tweetHandler.CreateTweet)
	group.GET("/tweets/:id", application.tweetHandler.GetTweet)
	group.GET("/tweets/users/:id", application.tweetHandler.GetUserTweets)
	group.DELETE("/tweets/:id", application.tweetHandler.DeleteTweet)
}
