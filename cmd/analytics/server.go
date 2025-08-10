package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func newServer() *Server {
	router := gin.Default()

	server := &Server{
		router: router,
	}

	return server
}

func (server *Server) Start() {
	log.Println("Starting service on :8080")
	if err := server.router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
