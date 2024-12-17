// controllers/unitOfMeasureController.go
package controllers

import (
	"farmers_market_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUnitOfMeasure creates a new unit of measure
func CreateUnitOfMeasure(c *gin.Context) {
	var unit models.UnitOfMeasure

	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Create(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit of measure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unit_of_measure": unit})
}

// GetUnitsOfMeasure retrieves all units of measure
func GetUnitsOfMeasure(c *gin.Context) {
	var units []models.UnitOfMeasure
	if err := db.Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve units of measure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"units_of_measure": units})
}

// UpdateUnitOfMeasure updates an existing unit of measure
func UpdateUnitOfMeasure(c *gin.Context) {
	unitID := c.Param("unitID")
	var updatedUnit models.UnitOfMeasure

	if err := c.ShouldBindJSON(&updatedUnit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var unit models.UnitOfMeasure
	if err := db.First(&unit, unitID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unit of measure not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve unit of measure"})
		}
		return
	}

	unit.Name = updatedUnit.Name
	unit.Abbreviation = updatedUnit.Abbreviation

	if err := db.Save(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update unit of measure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unit_of_measure": unit})
}

// DeleteUnitOfMeasure deletes a unit of measure by ID
func DeleteUnitOfMeasure(c *gin.Context) {
	unitID := c.Param("unitID")

	if err := db.Delete(&models.UnitOfMeasure{}, unitID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete unit of measure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unit of measure deleted successfully"})
}
