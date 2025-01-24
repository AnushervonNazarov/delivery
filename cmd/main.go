package main

import (
	"context"
	"delivery/configs"
	"delivery/internal/controllers"
	"delivery/logger"
	"delivery/pkg/db"
	"delivery/server"
	"fmt"
	"log"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("/home/hp/go/src/delivery/.env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	if err := configs.ReadSettings(); err != nil {
		log.Fatalf("Error reading settings: %s", err)
	}

	if err := logger.Init(); err != nil {
		log.Fatalf("Logger initialization error: %s", err)
	}

	if err := db.ConnectToDB(); err != nil {
		log.Fatalf("Database connection error: %s", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("Database migration error: %s", err)
	}

	mainServer := new(server.Server)
	go func() {
		if err := mainServer.Run(configs.AppSettings.AppParams.PortRun, controllers.RunRouts()); err != nil {
			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()

	// Waiting for signal to terminate
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Printf("\nStart of program completion\n")

	// Close the connection to the database if necessary
	if sqlDB, err := db.GetDBConn().DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Fatalf("Error closing connection to DB: %s", err)
		}
	} else {
		log.Fatalf("Error getting *sql.DB from GORM: %s", err)
	}
	fmt.Println("The connection to the database was closed successfully.")

	// Using a context with a timeout to shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mainServer.Shutdown(ctx); err != nil {
		log.Fatalf("Error while shutting down the server: %s", err)
	}

	fmt.Println("HTTP service successfully shut down")
	fmt.Println("End of program completion")
}
