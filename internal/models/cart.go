package models

type Cart struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	UserID   uint `gorm:"not null" json:"user_id"`
	ItemID   uint `gorm:"not null" json:"item_id"`
	Quantity int  `gorm:"not null" json:"quantity"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Item Item `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
}
