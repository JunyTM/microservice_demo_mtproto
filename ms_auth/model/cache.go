package model

type CacheMem struct {
	ID    int64  `json:"id" gorm:"primaryKey"`
	Email string `json:"email"`
}
