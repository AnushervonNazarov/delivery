package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null:unique"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SignInInput struct {
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"not null:unique"`
}
