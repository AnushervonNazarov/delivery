package db

import "blogging_platform/internal/models"

func Migrate() error {
	err := dbConn.AutoMigrate(
		models.User{})
	if err != nil {
		return err
	}
	return nil
}
