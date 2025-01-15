package repository

import (
	"delivery/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItems() ([]models.Item, error)
	GetItemByID(itemID int) (models.Item, error)
	UpdateItemStock(itemID int, quantity int) error
}

type itemRepository struct {
	db *gorm.DB
}

// NewItemRepository returns a new instance of ItemRepository using GORM DB connection
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

// GetAllItems retrieves all items from the database
func (r *itemRepository) GetAllItems() ([]models.Item, error) {
	var items []models.Item
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetItemByID retrieves an item by its ID
func (r *itemRepository) GetItemByID(itemID int) (models.Item, error) {
	var item models.Item
	// Log the itemID to ensure it's passed correctly
	fmt.Printf("Fetching item with ID: %d\n", itemID)

	err := r.db.First(&item, itemID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("Item with ID %d not found in the database\n", itemID)
			return models.Item{}, errors.New("item not found")
		}
		// Log any other database errors
		fmt.Printf("Error retrieving item with ID %d: %v\n", itemID, err)
		return models.Item{}, err
	}

	// Log the item details if found
	fmt.Printf("Item fetched: %+v\n", item)

	return item, nil
}

// UpdateItemStock updates the stock of an item
func (r *itemRepository) UpdateItemStock(itemID int, quantity int) error {
	// Here, `Updates` is used to decrement the stock based on the provided quantity
	if err := r.db.Model(&models.Item{}).Where("id = ?", itemID).UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
		return err
	}
	return nil
}
