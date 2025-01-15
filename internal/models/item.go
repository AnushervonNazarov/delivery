package models

type Item struct {
	ID    uint    `gorm:"primaryKey" json:"id"`
	Name  string  `gorm:"type:varchar(100);not null" json:"name"`
	Price float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock int     `gorm:"not null" json:"stock"`
}
