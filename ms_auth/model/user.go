package model

import (
	"time"

	"gorm.io/gorm"
)

// Model in MySQL
type User struct {
	ID        int64          `json:"id" gorm:"primary_key"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"UNIQUE;index"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// =================================================================
// Model message object from Client to Server
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	UserId        int64  `json:"user_id"`
	Email         string `json:"email"`
	ServerMessage string `json:"sv_message"`
}