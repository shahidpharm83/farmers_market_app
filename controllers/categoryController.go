// controllers/categoryController.go
package controllers

import (
	"farmers_market_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// GetCategories retrieves all categories
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// UpdateCategory updates an existing category
func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("categoryID")
	var updatedCategory models.Category

	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var category models.Category
	if err := db.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve category"})
		}
		return
	}

	category.Name = updatedCategory.Name
	category.Description = updatedCategory.Description

	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// DeleteCategory deletes a category by ID
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("categoryID")

	if err := db.Delete(&models.Category{}, categoryID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
