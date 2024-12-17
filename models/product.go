// models/product.go
package models

import (
	"time"
)

type Product struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Price             float64    `json:"price"`
	Discount          float64    `json:"discount"`           // New field for discount percentage
	IsPromoSale       bool       `json:"is_promo_sale"`      // New field for promotional sale
	SeasonExpiryDate  *time.Time `json:"season_expiry_date"` // New field for seasonal expiry date
	Stock             int        `json:"stock"`
	CategoryID        uint       `json:"category_id" gorm:"not null"`        // Foreign Key
	UnitOfMeasureID   uint       `json:"unit_of_measure_id" gorm:"not null"` // Foreign Key
	ImageURL          string     `json:"image_url"`
	VideoURL          string     `json:"video_url"`
	Min_order_qty     float64    `json:"min_order_qty"`
	SellerID          uint       `json:"seller_id" gorm:"not null"` // Foreign Key from User
	Rating            float64    `json:"rating" gorm:"default:0"`
	DeliveryTime      int        `json:"delivery_time"`       // Duration in minutes or seconds
	DeliveryTimeRules string     `json:"delivery_time_rules"` // New field for rules regarding delivery time
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relationships
	Category      Category      `json:"category" gorm:"foreignKey:CategoryID"`
	UnitOfMeasure UnitOfMeasure `json:"unit_of_measure" gorm:"foreignKey:UnitOfMeasureID"`
	Seller        User          `json:"seller" gorm:"foreignKey:SellerID"` // Relationship to User
}
