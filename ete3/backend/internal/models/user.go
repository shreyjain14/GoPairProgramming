package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" binding:"required,unique"`
	Email    string `json:"email" binding:"required,email,unique"`
	Password string `json:"-" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
