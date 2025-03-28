package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserId      int       `json:"userid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
