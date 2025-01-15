package models

type OrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	OrderID  uint    `gorm:"not null" json:"order_id"`
	ItemID   uint    `gorm:"not null" json:"item_id"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	// Relationships
	Order Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Item  Item  `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
}
