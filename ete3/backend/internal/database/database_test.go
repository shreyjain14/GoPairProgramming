package database

import (
	"ete3/internal/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set up test database
	os.Setenv("DB_PATH", ":memory:")
	InitDB()
	
	// Run tests
	code := m.Run()
	
	// Clean up
	os.Exit(code)
}

func TestCreateMovie(t *testing.T) {
	movie := &models.Movie{
		Title:       "Test Movie",
		Description: "Test Description",
		Duration:    120,
	}

	err := CreateMovie(movie)
	assert.NoError(t, err)
	assert.NotZero(t, movie.ID)

	// Verify movie was created
	createdMovie, err := GetMovieByID(movie.ID)
	assert.NoError(t, err)
	assert.Equal(t, movie.Title, createdMovie.Title)
	assert.Equal(t, movie.Description, createdMovie.Description)
	assert.Equal(t, movie.Duration, createdMovie.Duration)
}

func TestGetMovies(t *testing.T) {
	// Create test movies
	movie1 := &models.Movie{
		Title:       "Movie 1",
		Description: "Description 1",
		Duration:    120,
	}
	movie2 := &models.Movie{
		Title:       "Movie 2",
		Description: "Description 2",
		Duration:    150,
	}

	CreateMovie(movie1)
	CreateMovie(movie2)

	// Get all movies
	movies, err := GetMovies()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(movies), 2)
}

func TestCreateBooking(t *testing.T) {
	// Create a test show and seats first
	showID := int64(1)
	seatIDs := []int64{1, 2, 3}

	response, err := CreateBooking(showID, seatIDs)
	assert.NoError(t, err)
	assert.NotZero(t, response.BookingID)
	assert.Equal(t, len(seatIDs), len(response.Seats))
}

func TestGetBookings(t *testing.T) {
	bookings, err := GetBookings()
	assert.NoError(t, err)
	assert.NotNil(t, bookings)
}

func TestCancelBooking(t *testing.T) {
	// Create a test booking first
	showID := int64(1)
	seatIDs := []int64{1, 2, 3}
	booking, err := CreateBooking(showID, seatIDs)
	assert.NoError(t, err)

	// Cancel the booking
	err = CancelBooking(booking.BookingID)
	assert.NoError(t, err)

	// Verify booking is cancelled
	bookings, err := GetBookings()
	assert.NoError(t, err)
	for _, b := range bookings {
		if b.ID == booking.BookingID {
			assert.True(t, b.Cancelled)
		}
	}
}

func TestGetTheaterLayout(t *testing.T) {
	showID := int64(1)
	layout, err := GetTheaterLayout(showID)
	assert.NoError(t, err)
	assert.NotNil(t, layout)
}

func TestGetAvailableSeats(t *testing.T) {
	showID := int64(1)
	seats, err := GetAvailableSeats(showID)
	assert.NoError(t, err)
	assert.NotNil(t, seats)
}

func TestGetShowsByMovie(t *testing.T) {
	movieID := int64(1)
	shows, err := GetShowsByMovie(movieID)
	assert.NoError(t, err)
	assert.NotNil(t, shows)
} 