# User Authentication API

A simple REST API built with Go, featuring user registration and authentication with SQLite3 database.

## Features

- User registration with unique username and email
- User login with JWT authentication
- SQLite3 database for data persistence
- ReDoc API documentation
- Password hashing using bcrypt

## Prerequisites

- Go 1.21 or higher
- SQLite3

## Setup

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## API Documentation

Access the API documentation at `http://localhost:8080/docs`

## API Endpoints

### Register User
- **POST** `/api/auth/register`
- Request body:
  ```json
  {
    "username": "unique_username",
    "email": "user@example.com",
    "password": "password123"
  }
  ```

### Login User
- **POST** `/api/auth/login`
- Request body:
  ```json
  {
    "username": "unique_username",
    "password": "password123"
  }
  ```

## Security

- Passwords are hashed using bcrypt
- JWT tokens are used for authentication
- Input validation is performed on all requests 