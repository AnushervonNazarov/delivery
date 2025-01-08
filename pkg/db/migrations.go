package db

import "delivery/internal/models"

func Migrate() error {
	err := dbConn.AutoMigrate(
		models.User{})
	if err != nil {
		return err
	}
	return nil
}
