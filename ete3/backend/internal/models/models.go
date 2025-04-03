package models

import (
	"encoding/json"
	"strings"
	"time"
)

type Movie struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"` // in minutes
	Genre       string    `json:"genre"`
	PosterURL   string    `json:"poster_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Theater struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Show struct {
	ID        int64     `json:"id"`
	MovieID   int64     `json:"movie_id"`
	TheaterID int64     `json:"theater_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Seat struct {
	ID         int64     `json:"id"`
	TheaterID  int64     `json:"theater_id"`
	RowNumber  int       `json:"row_number"`
	SeatNumber int       `json:"seat_number"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Booking struct {
	ID        int64     `json:"id"`
	ShowID    int64     `json:"show_id"`
	SeatID    int64     `json:"seat_id"`
	Status    string    `json:"status"` // "pending", "confirmed", "cancelled"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookingRequest struct {
	ShowID  int64   `json:"show_id" binding:"required"`
	SeatIDs []int64 `json:"seat_ids" binding:"required"`
}

type BookingResponse struct {
	BookingID int64  `json:"booking_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// TheaterLayout represents a visual layout of seats in a theater
type TheaterLayout struct {
	TheaterID int64          `json:"theater_id"`
	Name      string         `json:"name"`
	Rows      int            `json:"rows"`
	Columns   int            `json:"columns"`
	Layout    [][]SeatStatus `json:"-"`      // Won't be directly marshalled
	LayoutMap string         `json:"layout"` // Custom marshalled field
}

// SeatStatus represents the status of a seat
type SeatStatus struct {
	ID     int64  `json:"id"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Status string `json:"status"` // "available", "booked", "selected"
}

// Custom marshalling for TheaterLayout
func (t TheaterLayout) MarshalJSON() ([]byte, error) {
	type Alias TheaterLayout

	// Convert the 2D layout to a compact string representation
	layoutMap := ""
	for _, row := range t.Layout {
		for _, seat := range row {
			switch seat.Status {
			case "available":
				layoutMap += "A"
			case "booked":
				layoutMap += "B"
			case "selected":
				layoutMap += "S"
			default:
				layoutMap += "X" // unavailable
			}
		}
		layoutMap += "|" // row separator
	}

	return json.Marshal(&struct {
		Alias
		LayoutMap string `json:"layout"`
	}{
		Alias:     Alias(t),
		LayoutMap: layoutMap,
	})
}

// Custom unmarshalling for TheaterLayout
func (t *TheaterLayout) UnmarshalJSON(data []byte) error {
	type Alias TheaterLayout
	aux := &struct {
		LayoutMap string `json:"layout"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert string representation back to 2D layout
	rows := strings.Split(aux.LayoutMap, "|")
	if len(rows) > 0 && len(rows[len(rows)-1]) == 0 {
		// Remove last empty element after final |
		rows = rows[:len(rows)-1]
	}

	t.Rows = len(rows)
	if t.Rows > 0 {
		t.Columns = len(rows[0])
	}

	// Initialize the layout array
	t.Layout = make([][]SeatStatus, t.Rows)
	for i := range t.Layout {
		t.Layout[i] = make([]SeatStatus, t.Columns)
	}

	// Fill in the layout data
	for i, row := range rows {
		for j, char := range row {
			if j >= t.Columns {
				continue
			}

			status := "unavailable"
			switch char {
			case 'A':
				status = "available"
			case 'B':
				status = "booked"
			case 'S':
				status = "selected"
			}

			t.Layout[i][j] = SeatStatus{
				Row:    i + 1,
				Column: j + 1,
				Status: status,
			}
		}
	}

	return nil
}
