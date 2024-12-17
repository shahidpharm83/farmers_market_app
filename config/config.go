package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase connects to the database and returns the DB instance.
func ConnectDatabase() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/farmers_market_db?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database!")
	}

	DB = database
	log.Println("Database connected successfully!")
	return database // Return the connected database instance
}

// LoadEnv loads environment variables from a .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Loading environment variables from system.")
	}

	// Check if JWT_SECRET_KEY environment variable is set
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("Error: JWT_SECRET environment variable not set. Exiting with status 1.")
		os.Exit(1)
	}
}

// GetEnv retrieves an environment variable or returns a fallback value
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: %s environment variable not set. Using fallback value: %s", key, fallback)
		return fallback
	}
	return value
}
