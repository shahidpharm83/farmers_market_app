package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase connects to the database and returns the DB instance.
func ConnectDatabase() *gorm.DB {
	// PostgreSQL DSN format:
	// "host=<HOST> user=<USERNAME> password=<PASSWORD> dbname=<DBNAME> port=<PORT> sslmode=<SSLMODE>"
	dsn := "host=dpg-ctgiqmd2ng1s738j303g-a user=admin dbname=farmers_market_db sslmode=disable password=0hj6px0nhvaahHftsmOFakdtYWbdN1AU"
	database, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal("Failed to connect to the PostgreSQL database:", err)
	}

	DB = database
	log.Println("PostgreSQL database connected successfully!")
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
