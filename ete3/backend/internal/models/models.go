package models

import "time"

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
