package database

import (
	"database/sql"
	"ete3/internal/models"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(user *models.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err = DB.Exec(query, user.Username, user.Email, string(hashedPassword))
	return err
}

func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password FROM users WHERE username = ?`
	err := DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
