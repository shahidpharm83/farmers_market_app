package controllers

import (
	"farmers_market_backend/config"
	"farmers_market_backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PlaceOrder creates a new order for a product
func PlaceOrder(c *gin.Context) {
	// CreateOrder creates a new order
	var order models.Order

	// Bind JSON input to order struct
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate product and users (buyer and seller)
	if err := validateOrderDependencies(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the product to calculate delivery date
	var product models.Product
	if err := db.First(&product, order.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Set OrderDateTime to current time
	order.OrderDateTime = time.Now()

	// Calculate DeliveryDateTime based on product's delivery rules
	order.DeliveryDateTime = order.OrderDateTime.Add(time.Duration(product.DeliveryTime) * time.Hour)

	// Create the order
	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// validateOrderDependencies checks if the product and users are valid
func validateOrderDependencies(order *models.Order) error {
	var product models.Product
	if err := db.First(&product, order.ProductID).Error; err != nil {
		return fmt.Errorf("invalid product ID")
	}

	var buyer models.User
	if err := db.First(&buyer, order.BuyerID).Error; err != nil {
		return fmt.Errorf("invalid buyer ID")
	}

	var seller models.User
	if err := db.First(&seller, order.SellerID).Error; err != nil {
		return fmt.Errorf("invalid seller ID")
	}

	return nil
}

// GetOrderDetails retrieves details of a specific order by ID
func GetOrderDetails(c *gin.Context) {
	var order models.Order
	if err := config.DB.Preload("Product").First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// UpdateOrder updates an existing order by ID
func UpdateOrder(c *gin.Context) {
	orderID := c.Param("orderID")
	var updatedOrder models.Order

	// Bind JSON data to the updatedOrder struct
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the existing order
	var order models.Order
	if err := db.First(&order, orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
		}
		return
	}

	// Update the order fields
	order.Status = updatedOrder.Status // Assuming there is a Status field
	order.TotalPrice = updatedOrder.TotalPrice
	// Add other fields to update as needed

	// Save the updated order
	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// DeleteOrder deletes an order by ID
func DeleteOrder(c *gin.Context) {
	orderID := c.Param("orderID")

	// Delete the order by ID
	if err := db.Delete(&models.Order{}, orderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
