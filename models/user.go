package models

import (
	"time"
)

// User struct
type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name"`
	Email           string    `json:"email" gorm:"unique"`
	Username        string    `json:"username" gorm:"unique"`
	Password        string    `json:"-"`
	ImageURL        string    `json:"image_url"`                             // Field for user image URL
	IsSeller        bool      `json:"is_seller"`                             // New field to indicate if the user is a seller
	DistrictID      *uint     `json:"district_id" gorm:"index"`              // Foreign key to District (pointer type)
	DeliveryAddress string    `json:"delivery_address"`                      // New field for delivery address
	MobileNumber    string    `json:"mobile_number"`                         // New field for mobile number
	District        District  `json:"district" gorm:"foreignKey:DistrictID"` // Relationship with District
	IsAdmin         bool      `json:"is_admin"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
