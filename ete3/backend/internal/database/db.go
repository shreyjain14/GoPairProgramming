package database

import (
	"database/sql"
	"ete3/internal/models"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	DB *sql.DB
	// Mutex for handling concurrent bookings
	bookingMutex sync.Mutex
)

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./cinema.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
	migrateDatabase()
}

// migrateDatabase runs any required database migrations
func migrateDatabase() {
	// Check if poster_url column exists in movies table
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('movies') WHERE name='poster_url'").Scan(&count)
	if err != nil {
		log.Printf("Error checking for poster_url column: %v", err)
		return
	}

	// If poster_url column doesn't exist, add it
	if count == 0 {
		log.Println("Adding poster_url column to movies table")
		_, err := DB.Exec("ALTER TABLE movies ADD COLUMN poster_url TEXT")
		if err != nil {
			log.Printf("Error adding poster_url column: %v", err)
			return
		}
		log.Println("poster_url column added successfully")
	}
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS movies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			duration INTEGER NOT NULL,
			genre TEXT,
			poster_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS theaters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			capacity INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS shows (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			movie_id INTEGER NOT NULL,
			theater_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			price REAL NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (movie_id) REFERENCES movies(id),
			FOREIGN KEY (theater_id) REFERENCES theaters(id)
		);`,
		`CREATE TABLE IF NOT EXISTS seats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			theater_id INTEGER NOT NULL,
			row_number INTEGER NOT NULL,
			seat_number INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (theater_id) REFERENCES theaters(id)
		);`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			show_id INTEGER NOT NULL,
			seat_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (show_id) REFERENCES shows(id),
			FOREIGN KEY (seat_id) REFERENCES seats(id)
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Movie operations
func GetMovies() ([]models.Movie, error) {
	rows, err := DB.Query("SELECT id, title, description, duration, genre, poster_url, created_at, updated_at FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Duration, &m.Genre, &m.PosterURL, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

// Show operations
func GetShowsByMovie(movieID int64) ([]models.Show, error) {
	rows, err := DB.Query(`
		SELECT s.id, s.movie_id, s.theater_id, s.start_time, s.end_time, s.price, s.created_at, s.updated_at
		FROM shows s
		WHERE s.movie_id = ? AND s.start_time > datetime('now')
		ORDER BY s.start_time`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []models.Show
	for rows.Next() {
		var s models.Show
		err := rows.Scan(&s.ID, &s.MovieID, &s.TheaterID, &s.StartTime, &s.EndTime, &s.Price, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shows = append(shows, s)
	}
	return shows, nil
}

// Seat operations
func GetAvailableSeats(showID int64) ([]models.Seat, error) {
	rows, err := DB.Query(`
		SELECT s.id, s.theater_id, s.row_number, s.seat_number, s.created_at, s.updated_at
		FROM seats s
		JOIN shows sh ON s.theater_id = sh.theater_id
		WHERE sh.id = ? AND s.id NOT IN (
			SELECT seat_id FROM bookings WHERE show_id = ? AND status != 'cancelled'
		)`, showID, showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var s models.Seat
		err := rows.Scan(&s.ID, &s.TheaterID, &s.RowNumber, &s.SeatNumber, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

// Booking operations with concurrency control
func CreateBooking(showID int64, seatIDs []int64) (*models.BookingResponse, error) {
	bookingMutex.Lock()
	defer bookingMutex.Unlock()

	// Start transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if seats are available
	for _, seatID := range seatIDs {
		var count int
		err := tx.QueryRow(`
			SELECT COUNT(*) FROM bookings 
			WHERE show_id = ? AND seat_id = ? AND status != 'cancelled'`, showID, seatID).Scan(&count)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return &models.BookingResponse{
				Status:  "failed",
				Message: "One or more seats are already booked",
			}, nil
		}
	}

	// Create bookings
	for _, seatID := range seatIDs {
		_, err := tx.Exec(`
			INSERT INTO bookings (show_id, seat_id, status)
			VALUES (?, ?, 'confirmed')`, showID, seatID)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.BookingResponse{
		Status:  "success",
		Message: "Booking confirmed successfully",
	}, nil
}

// Get all bookings
func GetBookings() ([]models.Booking, error) {
	rows, err := DB.Query(`
		SELECT id, show_id, seat_id, status, created_at, updated_at
		FROM bookings
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		err := rows.Scan(&b.ID, &b.ShowID, &b.SeatID, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

// Cancel booking
func CancelBooking(bookingID int64) error {
	result, err := DB.Exec(`
		UPDATE bookings 
		SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 'confirmed'`, bookingID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CreateMovie adds a new movie to the database and creates shows with seats
func CreateMovie(movie *models.Movie) error {
	// Start transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert movie
	result, err := tx.Exec(`
		INSERT INTO movies (title, description, duration, genre, poster_url)
		VALUES (?, ?, ?, ?, ?)`,
		movie.Title, movie.Description, movie.Duration, movie.Genre, movie.PosterURL)

	if err != nil {
		return err
	}

	movieID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	movie.ID = movieID

	// Create a theater if none exists
	var theaterID int64
	err = tx.QueryRow("SELECT id FROM theaters LIMIT 1").Scan(&theaterID)
	if err != nil {
		// Create a theater if none exists
		res, err := tx.Exec(`
			INSERT INTO theaters (name, capacity)
			VALUES (?, ?)`,
			"Main Theater", 140)
		if err != nil {
			return err
		}
		theaterID, err = res.LastInsertId()
		if err != nil {
			return err
		}

		// Create seats for the theater (A1-A20 to G1-G20)
		for row := 0; row < 7; row++ {
			for seat := 1; seat <= 20; seat++ {
				_, err = tx.Exec(`
					INSERT INTO seats (theater_id, row_number, seat_number)
					VALUES (?, ?, ?)`,
					theaterID, row+1, seat)
				if err != nil {
					return err
				}
			}
		}
	}

	// Create shows for the next 7 days at 6pm, 8pm, and 10pm
	for day := 0; day < 7; day++ {
		// 6:00 PM show
		_, err = tx.Exec(`
			INSERT INTO shows (movie_id, theater_id, start_time, end_time, price)
			VALUES (?, ?, datetime('now', '+' || ? || ' days', '18:00:00'), 
				datetime('now', '+' || ? || ' days', '18:00:00', '+' || ? || ' minutes'), ?)`,
			movieID, theaterID, day, day, movie.Duration, 10.00)
		if err != nil {
			return err
		}

		// 8:00 PM show
		_, err = tx.Exec(`
			INSERT INTO shows (movie_id, theater_id, start_time, end_time, price)
			VALUES (?, ?, datetime('now', '+' || ? || ' days', '20:00:00'), 
				datetime('now', '+' || ? || ' days', '20:00:00', '+' || ? || ' minutes'), ?)`,
			movieID, theaterID, day, day, movie.Duration, 12.00)
		if err != nil {
			return err
		}

		// 10:00 PM show
		_, err = tx.Exec(`
			INSERT INTO shows (movie_id, theater_id, start_time, end_time, price)
			VALUES (?, ?, datetime('now', '+' || ? || ' days', '22:00:00'), 
				datetime('now', '+' || ? || ' days', '22:00:00', '+' || ? || ' minutes'), ?)`,
			movieID, theaterID, day, day, movie.Duration, 8.00)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateMovie updates an existing movie in the database
func UpdateMovie(movie *models.Movie) error {
	_, err := DB.Exec(`
		UPDATE movies 
		SET title = ?, description = ?, duration = ?, genre = ?, poster_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		movie.Title, movie.Description, movie.Duration, movie.Genre, movie.PosterURL, movie.ID)

	return err
}

// GetMovieByID retrieves a movie by its ID
func GetMovieByID(movieID int64) (*models.Movie, error) {
	movie := &models.Movie{}
	err := DB.QueryRow(`
		SELECT id, title, description, duration, genre, poster_url, created_at, updated_at 
		FROM movies 
		WHERE id = ?`, movieID).Scan(
		&movie.ID, &movie.Title, &movie.Description, &movie.Duration,
		&movie.Genre, &movie.PosterURL, &movie.CreatedAt, &movie.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

// User operations
func CreateUser(req *models.RegisterRequest) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		INSERT INTO users (username, email, password)
		VALUES (?, ?, ?)`,
		req.Username, req.Email, string(hashedPassword))

	return err
}

func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow(`
		SELECT id, username, email, password 
		FROM users 
		WHERE username = ?`, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}
