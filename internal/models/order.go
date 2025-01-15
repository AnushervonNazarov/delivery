package models

import "time"

type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	TotalAmount float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	// Relationships
	User       User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}
