package services

import (
	"farmers_market_backend/config"
	"farmers_market_backend/models"

	"gorm.io/gorm"
)

func TopUpWallet(userID uint, amount float64) error {
	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			wallet = models.Wallet{UserID: userID, Balance: amount}
			return config.DB.Create(&wallet).Error
		}
		return err
	}

	wallet.Balance += amount
	return config.DB.Save(&wallet).Error
}
