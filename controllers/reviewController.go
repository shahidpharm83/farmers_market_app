// controllers/reviewController.go
package controllers

import (
	"farmers_market_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateReview handles creating a new review
func CreateReview(c *gin.Context) {
	var review models.Review

	// Bind JSON to the review struct
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Save the review in the database
	if err := db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review created successfully", "review": review})
}

// GetReviews retrieves all reviews for a specific product
func GetReviews(c *gin.Context) {
	var reviews []models.Review
	productID := c.Param("productID")

	// Find all reviews for a specific product
	if err := db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// UpdateReview handles updating an existing review
func UpdateReview(c *gin.Context) {
	var review models.Review
	reviewID := c.Param("reviewID")

	// Find the review by ID
	if err := db.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Bind JSON to the review struct
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Update the review in the database
	if err := db.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully", "review": review})
}

// DeleteReview handles deleting a review
func DeleteReview(c *gin.Context) {
	reviewID := c.Param("reviewID")

	// Delete the review by ID
	if err := db.Delete(&models.Review{}, reviewID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
