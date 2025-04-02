package main

import (
	"ete3/internal/database"
	"ete3/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	// Create Gin router
	r := gin.Default()

	// Serve static files from docs directory
	r.Static("/docs", "./docs")
	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
