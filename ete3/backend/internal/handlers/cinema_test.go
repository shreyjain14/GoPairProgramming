package handlers

import (
	"bytes"
	"encoding/json"
	"ete3/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestGetMovies(t *testing.T) {
	router := setupRouter()
	router.GET("/api/cinema/movies", GetMovies)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cinema/movies", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMovie(t *testing.T) {
	router := setupRouter()
	router.GET("/api/cinema/movies/:id", GetMovie)

	tests := []struct {
		name       string
		movieID    string
		wantStatus int
	}{
		{
			name:       "Valid Movie ID",
			movieID:    "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Movie ID",
			movieID:    "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-existent Movie ID",
			movieID:    "999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/cinema/movies/"+tt.movieID, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCreateMovie(t *testing.T) {
	router := setupRouter()
	router.POST("/api/cinema/movies", CreateMovie)

	tests := []struct {
		name       string
		movie      models.Movie
		wantStatus int
	}{
		{
			name: "Valid Movie",
			movie: models.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Duration:    120,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid Movie - Missing Title",
			movie: models.Movie{
				Description: "Test Description",
				Duration:    120,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.movie)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/cinema/movies", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCreateBooking(t *testing.T) {
	router := setupRouter()
	router.POST("/api/cinema/bookings", CreateBooking)

	tests := []struct {
		name       string
		booking    models.BookingRequest
		wantStatus int
	}{
		{
			name: "Valid Booking",
			booking: models.BookingRequest{
				ShowID:  1,
				SeatIDs: []int64{1, 2, 3},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid Booking - No Seats",
			booking: models.BookingRequest{
				ShowID:  1,
				SeatIDs: []int64{},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.booking)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/cinema/bookings", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCancelBooking(t *testing.T) {
	router := setupRouter()
	router.DELETE("/api/cinema/bookings/:id", CancelBooking)

	tests := []struct {
		name       string
		bookingID  string
		wantStatus int
	}{
		{
			name:       "Valid Booking ID",
			bookingID:  "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Booking ID",
			bookingID:  "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-existent Booking ID",
			bookingID:  "999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/cinema/bookings/"+tt.bookingID, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetTheaterLayout(t *testing.T) {
	router := setupRouter()
	router.GET("/api/cinema/shows/:id/layout", GetTheaterLayout)

	tests := []struct {
		name       string
		showID     string
		wantStatus int
	}{
		{
			name:       "Valid Show ID",
			showID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Show ID",
			showID:     "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-existent Show ID",
			showID:     "999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/cinema/shows/"+tt.showID+"/layout", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
} 