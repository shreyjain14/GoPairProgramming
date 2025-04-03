package main

import (
	"ete3/internal/database"
	"ete3/internal/handlers"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting cinema booking application...")

	// Initialize database
	fmt.Println("Initializing database...")
	database.InitDB()
	fmt.Println("Database initialized successfully")

	// Create Gin router
	fmt.Println("Setting up Gin router...")
	r := gin.Default()

	// Configure CORS middleware
	fmt.Println("Configuring CORS to allow all origins...")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}))

	// Serve static files from docs directory
	r.Static("/docs", "./docs")
	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")

	// API routes
	api := r.Group("/api")
	{
		// Cinema routes
		cinema := api.Group("/cinema")
		{
			// Movies
			cinema.GET("/movies", handlers.GetMovies)
			cinema.GET("/movies/:id", handlers.GetMovie)
			cinema.POST("/movies", handlers.CreateMovie)
			cinema.PUT("/movies/:id", handlers.UpdateMovie)
			cinema.GET("/movies/:id/shows", handlers.GetShowsByMovie)

			// Shows and Seats
			cinema.GET("/shows/:id/seats", handlers.GetAvailableSeats)
			cinema.GET("/shows/:id/layout", handlers.GetTheaterLayout)

			// Bookings
			bookings := cinema.Group("/bookings")
			{
				bookings.POST("", handlers.CreateBooking)
				bookings.GET("", handlers.GetBookings)
				bookings.DELETE("/:id", handlers.CancelBooking)
			}
		}
	}

	// Start server
	fmt.Println("Starting server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
