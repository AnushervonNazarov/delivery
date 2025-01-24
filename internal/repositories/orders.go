package repository

import (
	"delivery/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	SaveOrder(order models.Order, items []models.Cart) (int, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) SaveOrder(order models.Order, items []models.Cart) (int, error) {
	// Start a transaction
	tx := r.db.Begin()

	// If there's an error starting the transaction, return it
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Insert the order into the orders table
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert each item in the order_items table
	for _, item := range items {
		// Fetch the item details to get the correct price
		var fetchedItem models.Item
		if err := tx.First(&fetchedItem, item.ItemID).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		orderItem := models.OrderItem{
			OrderID:  order.ID,
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
			Price:    float64(item.Quantity) * fetchedItem.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Commit the transaction if everything goes well
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return int(order.ID), nil
}
