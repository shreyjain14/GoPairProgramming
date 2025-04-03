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
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
