// models/order.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ProductID        uint      `json:"product_id" gorm:"not null"`         // Foreign Key from Product
	BuyerID          uint      `json:"buyer_id" gorm:"not null"`           // Foreign Key from User (Buyer)
	SellerID         uint      `json:"seller_id" gorm:"not null"`          // Foreign Key from User (Seller)
	Quantity         int       `json:"quantity" gorm:"not null"`           // Quantity of product ordered
	TotalPrice       float64   `json:"total_price" gorm:"not null"`        // Total price of the order
	Status           string    `json:"status" gorm:"default:'Pending'"`    // Status of the order
	OrderDateTime    time.Time `json:"order_date_time" gorm:"not null"`    // Date and time when the order was placed
	DeliveryDateTime time.Time `json:"delivery_date_time" gorm:"not null"` // Calculated delivery date and time

	// Relationships
	Product Product `json:"product" gorm:"foreignKey:ProductID"` // Relationship to Product
	Buyer   User    `json:"buyer" gorm:"foreignKey:BuyerID"`     // Relationship to User (Buyer)
	Seller  User    `json:"seller" gorm:"foreignKey:SellerID"`   // Relationship to User (Seller)
}
