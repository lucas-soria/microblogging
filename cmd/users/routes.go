package main

import (
	"github.com/lucas-soria/microblogging/cmd/users/middleware"

	"github.com/gin-gonic/gin"
)

// setupRoutes configures all the routes for the application
func addRoutes(server *Server, application *Application) {
	healthCheck(server.router)
	groupV1 := newGroup(server.router, "/v1")
	usersRoutes(groupV1, application)
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

func usersRoutes(group *gin.RouterGroup, application *Application) {
	group.POST("/users", application.userHandler.CreateUser)
	group.GET("/users/:id", application.userHandler.GetUser)
	group.GET("/users/:id/followers", application.userHandler.GetUserFollowers)
	group.GET("/users/:id/followees", application.userHandler.GetUserFollowees)

	protectedGroup := group.Group("")
	protectedGroup.Use(middleware.AuthMiddleware())
	protectedGroup.DELETE("/users/:id", application.userHandler.DeleteUser)
	protectedGroup.POST("/users/:id/follow", application.userHandler.FollowUser)
	protectedGroup.POST("/users/:id/unfollow", application.userHandler.UnfollowUser)
}
