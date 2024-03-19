package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `json:"id" gorm:"primary_key"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
