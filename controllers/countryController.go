package controllers

import (
	"farmers_market_backend/config"
	"farmers_market_backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCountries retrieves all countries
func GetCountries(c *gin.Context) {
	var countries []models.Country

	// Ensure db is initialized
	if config.DB == nil {
		log.Println("Database connection is nil")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	// Initialize countries slice to avoid nil pointer dereference
	if len(countries) == 0 {
		countries = []models.Country{}
	}

	// Fetch countries from the database
	if err := config.DB.Find(&countries).Error; err != nil {
		log.Printf("Error fetching countries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching countries"})
		return
	}

	// Return the list of countries
	c.JSON(http.StatusOK, countries)
}

// GetCountry retrieves a country by ID
func GetCountry(c *gin.Context) {
	var country models.Country
	id := c.Param("id")
	if err := config.DB.First(&country, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}
	c.JSON(http.StatusOK, country)
}

// CreateCountry creates a new country
func CreateCountry(c *gin.Context) {
	var country models.Country

	// Bind JSON body with error handling
	if err := c.ShouldBindJSON(&country); err != nil {
		log.Printf("JSON Binding Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Attempt to create the country in the database with detailed error logging
	if err := config.DB.Create(&country).Error; err != nil {
		log.Printf("Database Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create country in the database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": country})
}

// UpdateCountry updates a country
func UpdateCountry(c *gin.Context) {
	var country models.Country
	id := c.Param("id")
	if err := config.DB.First(&country, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}
	if err := c.ShouldBindJSON(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	db.Save(&country)
	c.JSON(http.StatusOK, country)
}

// DeleteCountry deletes a country
func DeleteCountry(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Country{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
