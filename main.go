package main

import (
	"farmers_market_backend/config"
	"farmers_market_backend/models"
	"farmers_market_backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// main is the entry point of the program. It initializes the Gin Router, loads
// environment variables and database configuration, initializes the database
// with migrations, sets up routes, and starts the server.
func main() {
	// Initialize Gin Router
	router := gin.Default()

	// Load environment variables and database configuration
	config.LoadEnv()
	database := config.ConnectDatabase() // Store the returned database instance

	// Initialize the database with migrations
	initDatabase(database)

	// Set up routes
	routes.InitializeRoutes(router)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// initDatabase performs database migrations for all models
func initDatabase(db *gorm.DB) {
	// Automatically migrate the schema
	if err := db.AutoMigrate(
		&models.Message{},
		&models.Review{},
		&models.Product{},
		&models.User{},
		&models.Country{},
		&models.Category{},
		&models.District{},
		&models.Order{},
		&models.Wallet{},
		&models.UnitOfMeasure{},
	); err != nil {
		log.Fatalf("Error during database migration: %v", err)
	}
}
