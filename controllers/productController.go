// controllers/productController.go
package controllers

import (
	"farmers_market_backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Bind JSON input to product struct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Create the product
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// GetProduct retrieves a product by ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := db.Preload("Category").Preload("UnitOfMeasure").Preload("Seller").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// GetAllProducts retrieves all products
func GetAllProducts(c *gin.Context) {
	var products []models.Product

	if err := db.Preload("Category").Preload("UnitOfMeasure").Preload("Seller").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	productID := c.Param("productID")
	var updatedProduct models.Product

	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the product exists
	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		}
		return
	}

	// Validate seller, category, and unit of measure
	if err := validateProductDependencies(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update product fields
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	product.Discount = updatedProduct.Discount
	product.IsPromoSale = updatedProduct.IsPromoSale
	product.SeasonExpiryDate = updatedProduct.SeasonExpiryDate
	product.Stock = updatedProduct.Stock
	product.CategoryID = updatedProduct.CategoryID
	product.UnitOfMeasureID = updatedProduct.UnitOfMeasureID
	product.ImageURL = updatedProduct.ImageURL
	product.VideoURL = updatedProduct.VideoURL
	product.Min_order_qty = updatedProduct.Min_order_qty
	product.SellerID = updatedProduct.SellerID // Update seller ID if necessary
	product.DeliveryTime = updatedProduct.DeliveryTime
	product.DeliveryTimeRules = updatedProduct.DeliveryTimeRules

	if err := db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// DeleteProduct deletes a product by ID
func DeleteProduct(c *gin.Context) {
	productID := c.Param("productID")

	if err := db.Delete(&models.Product{}, productID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// validateProductDependencies checks if the related entities exist
func validateProductDependencies(product *models.Product) error {
	var seller models.User
	if err := db.First(&seller, product.SellerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("seller not found")
		}
		return err
	}

	var category models.Category
	if err := db.First(&category, product.CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("category not found")
		}
		return err
	}

	var unit models.UnitOfMeasure
	if err := db.First(&unit, product.UnitOfMeasureID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("unit of measure not found")
		}
		return err
	}

	return nil
}
