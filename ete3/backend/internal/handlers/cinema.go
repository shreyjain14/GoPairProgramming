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

// GetTheaterLayout returns the layout of a theater with seat status for a specific show
func GetTheaterLayout(c *gin.Context) {
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

	// Get show information to get theater ID
	show, err := database.GetShowByID(showID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch show"})
		return
	}

	// Get theater information
	theater, err := database.GetTheaterByID(show.TheaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch theater"})
		return
	}

	// Get all seats for the theater
	seats, err := database.GetAllSeatsForTheater(show.TheaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}

	// Get booked seats for the show
	bookedSeats, err := database.GetBookedSeatsForShow(showID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booked seats"})
		return
	}

	// Create a map of booked seat IDs for quick lookup
	bookedSeatMap := make(map[int64]bool)
	for _, seat := range bookedSeats {
		bookedSeatMap[seat.ID] = true
	}

	// Find max row and column to determine theater dimensions
	maxRow, maxCol := 0, 0
	for _, seat := range seats {
		if seat.RowNumber > maxRow {
			maxRow = seat.RowNumber
		}
		if seat.SeatNumber > maxCol {
			maxCol = seat.SeatNumber
		}
	}

	// Create theater layout
	layout := models.TheaterLayout{
		TheaterID: theater.ID,
		Name:      theater.Name,
		Rows:      maxRow,
		Columns:   maxCol,
		Layout:    make([][]models.SeatStatus, maxRow),
	}

	// Initialize the layout with all seats marked as unavailable
	for i := range layout.Layout {
		layout.Layout[i] = make([]models.SeatStatus, maxCol)
		for j := range layout.Layout[i] {
			layout.Layout[i][j] = models.SeatStatus{
				Row:    i + 1,
				Column: j + 1,
				Status: "unavailable",
			}
		}
	}

	// Update the layout with actual seats and their status
	for _, seat := range seats {
		row := seat.RowNumber - 1
		col := seat.SeatNumber - 1

		if row >= 0 && row < maxRow && col >= 0 && col < maxCol {
			status := "available"
			if bookedSeatMap[seat.ID] {
				status = "booked"
			}

			layout.Layout[row][col] = models.SeatStatus{
				ID:     seat.ID,
				Row:    seat.RowNumber,
				Column: seat.SeatNumber,
				Status: status,
			}
		}
	}

	c.JSON(http.StatusOK, layout)
}
