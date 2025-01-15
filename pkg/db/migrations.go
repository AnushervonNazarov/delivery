package db

import "delivery/internal/models"

func Migrate() error {
	err := dbConn.AutoMigrate(
		models.User{},
		models.Menu{},
		models.Order{},
		models.OrderItem{},
		models.Item{},
		models.Cart{})
	if err != nil {
		return err
	}
	return nil
}
