package controllers

import (
	"farmers_market_backend/config"
	"farmers_market_backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TopUpWallet adds funds to a user's wallet
func TopUpWallet(c *gin.Context) {
	var request struct {
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)
	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		wallet = models.Wallet{UserID: userID, Balance: request.Amount}
		if err := config.DB.Create(&wallet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
			return
		}
	} else {
		wallet.Balance += request.Amount
		if err := config.DB.Save(&wallet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up wallet"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet topped up successfully", "balance": wallet.Balance})
}

// GetWalletBalance retrieves the current balance of a user's wallet
func GetWalletBalance(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": wallet.Balance})
}
