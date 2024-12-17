package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	Balance float64 `json:"balance" gorm:"default:0"`
}
