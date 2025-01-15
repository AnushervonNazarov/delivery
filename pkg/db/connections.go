package db

import (
	"delivery/configs"
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func ConnectToDB() error {

	connStr := fmt.Sprintf(`host=%s 
									port=%s 
									user=%s 
									dbname=%s 
									password=%s`,
		configs.AppSettings.PostgresParams.Host,
		configs.AppSettings.PostgresParams.Port,
		configs.AppSettings.PostgresParams.User,
		configs.AppSettings.PostgresParams.Database,
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return err
	}

	fmt.Println("Connected to database")

	dbConn = db

	return nil

	// connStr := "user=postgres password=q123 dbname=postgres host=localhost port=5433 sslmode=disable"

	// db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Connected to database")

	// dbConn = db
	// return nil
}

func CloseDBConn() error {
	// err := dbConn.Close()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func GetDBConn() *gorm.DB {
	return dbConn
}
