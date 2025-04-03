package handlers

import (
	"database/sql"
	"ete3/internal/database"
	"ete3/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMovies returns all available movies
func GetMovies(c *gin.Context) {
	movies, err := database.GetMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// GetShowsByMovie returns all shows for a specific movie
func GetShowsByMovie(c *gin.Context) {
	movieIDStr := c.Param("id")
	if movieIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
		return
	}

	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	shows, err := database.GetShowsByMovie(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shows"})
		return
	}
	c.JSON(http.StatusOK, shows)
}

// GetAvailableSeats returns all available seats for a specific show
func GetAvailableSeats(c *gin.Context) {
	showIDStr := c.Param("id")
	if showIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Show ID is required"})
		return
	}

	showID, err := strconv.ParseInt(showIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid show ID"})
		return
	}

	seats, err := database.GetAvailableSeats(showID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available seats"})
		return
	}
	c.JSON(http.StatusOK, seats)
}

// CreateBooking creates a new booking
func CreateBooking(c *gin.Context) {
	var req models.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := database.CreateBooking(req.ShowID, req.SeatIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBookings returns all bookings
func GetBookings(c *gin.Context) {
	bookings, err := database.GetBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// CancelBooking cancels a specific booking
func CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	if bookingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID is required"})
		return
	}

	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	err = database.CancelBooking(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found or already cancelled"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

// CreateMovie creates a new movie
func CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.CreateMovie(&movie); err != nil {
		log.Printf("Error creating movie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// UpdateMovie updates an existing movie
func UpdateMovie(c *gin.Context) {
	movieIDStr := c.Param("id")
	if movieIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
		return
	}

	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie.ID = movieID
	if err := database.UpdateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// GetMovie retrieves a movie by ID
func GetMovie(c *gin.Context) {
	movieIDStr := c.Param("id")
	if movieIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
		return
	}

	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := database.GetMovieByID(movieID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie"})
		return
	}

	c.JSON(http.StatusOK, movie)
}
