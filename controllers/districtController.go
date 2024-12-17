package controllers

import (
	"farmers_market_backend/config"
	"farmers_market_backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDistricts retrieves all districts
func GetDistricts(c *gin.Context) {
	var districts []models.District
	if err := db.Find(&districts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve districts"})
		return
	}
	c.JSON(http.StatusOK, districts)
}

// GetDistrict retrieves a district by ID
func GetDistrict(c *gin.Context) {
	var district models.District
	id := c.Param("id")
	if err := db.First(&district, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "District not found"})
		return
	}
	c.JSON(http.StatusOK, district)
}

// CreateDistrict creates a new district
func CreateDistrict(c *gin.Context) {

	// Initialize the database connection here
	var district models.District

	// Attempt to bind JSON to the district struct
	if err := c.ShouldBindJSON(&district); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Check if db is nil before using it
	if config.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	var country models.Country
	if err := config.DB.First(&country, district.CountryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Country ID", "details": err.Error()})
		return
	}

	// Create the district
	if err := config.DB.Create(&district).Error; err != nil {
		log.Printf("Database Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create district", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": district})
}

// UpdateDistrict updates a district
func UpdateDistrict(c *gin.Context) {
	var district models.District
	id := c.Param("id")
	if err := db.First(&district, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "District not found"})
		return
	}
	if err := c.ShouldBindJSON(&district); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	db.Save(&district)
	c.JSON(http.StatusOK, district)
}

// DeleteDistrict deletes a district
func DeleteDistrict(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&models.District{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "District not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
