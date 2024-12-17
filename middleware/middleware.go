package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"farmers_market_backend/config"
	"farmers_market_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware checks if the user is logged in by verifying the JWT token
// JWTSecretKey is loaded from the environment variable
var JWTSecretKey = []byte(config.GetEnv("JWT_SECRET", ""))

// getJWTSecretKey loads the JWT secret from the environment or logs an error
// getJWTSecretKey gets the JWT secret key from the environment variable

// AuthMiddleware checks if the user is logged in by verifying the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header (or wherever your token is stored)
		token := c.GetHeader("Authorization")

		// Log the token for debugging purposes
		log.Printf("Received Token: %v", token)

		// Validate the token and extract the user ID
		userID, err := validateTokenAndGetUserID(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Log the user ID
		log.Printf("Extracted User ID: %v", userID)

		// Set the user ID in the context for further use in the handlers
		c.Set("id", userID)

		// Query the user from the database to get the role (IsAdmin)
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Log the IsAdmin flag
		log.Printf("User IsAdmin: %v", user.IsAdmin)

		// Set the user role (IsAdmin) in the context for access in other middleware or handlers
		c.Set("IsAdmin", user.IsAdmin)

		// Continue to the next handler
		c.Next()
	}
}

// validateTokenAndGetUserID validates the JWT token and extracts user ID
// validateTokenAndGetUserID validates the JWT token and extracts the user ID
// validateTokenAndGetUserID validates the JWT token and extracts the user ID
func validateTokenAndGetUserID(tokenString string) (uint, error) {
	if tokenString == "" {
		log.Printf("Error: Token is missing")
		return 0, fmt.Errorf("token missing")
	}

	// Trim the 'Bearer ' prefix from the token if it exists
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	log.Printf("Token after trimming 'Bearer ': %v", tokenString)

	log.Printf("Token before parsing: %v", tokenString)
	// Parse the JWT token
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	log.Printf("Validating token with secret: %v", JWTSecretKey)
	// 	return JWTSecretKey, nil
	// })

	// if err != nil {
	// 	log.Printf("Error parsing token: %v", err)
	// 	if strings.Contains(err.Error(), "signature is invalid") {
	// 		log.Printf("Invalid token signature")
	// 	} else if strings.Contains(err.Error(), "token is expired") {
	// 		log.Printf("Token is expired")
	// 	} else {
	// 		log.Printf("Error: %v", err)
	// 	}
	// 	return 0, fmt.Errorf("invalid or expired token")
	// }
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		log.Printf("Validating token with secret: %v", JWTSecretKey)
		return JWTSecretKey, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		if strings.Contains(err.Error(), "signature is invalid") {
			log.Printf("Invalid token signature: %v", err)
		} else if strings.Contains(err.Error(), "token is expired") {
			log.Printf("Token is expired")
		} else {
			log.Printf("Error: %v", err)
		}
		return 0, fmt.Errorf("invalid or expired token")
	}

	if !token.Valid {
		log.Printf("Invalid token signature")
		return 0, fmt.Errorf("invalid token")
	}

	// Extract claims and get user ID
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			log.Printf("Extracted User ID: %v", userID)
			return uint(userID), nil
		}
		log.Printf("Error: User ID not found in token claims")
		return 0, fmt.Errorf("user ID not found in token claims")
	}

	log.Printf("Error: Invalid token claims")
	return 0, fmt.Errorf("invalid token claims")
}

func IsSellerOfTheItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: Extract item ID from the URL parameter
		productID := c.Param("itemID")
		userID := c.MustGet("userID").(uint) // Assume userID is set in the context after successful login

		var product models.Product // Assume Product is your model
		if err := config.DB.Where("id = ? AND seller_id = ?", productID, userID).First(&product).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this item"})
			c.Abort()
			return
		}

		c.Next()
	}
}
func IsCustomerForThisOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("orderID")
		userID := c.MustGet("userID").(uint)

		var order models.Order // Assume Order is your model
		if err := config.DB.Where("id = ? AND customer_id = ?", orderID, userID).First(&order).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this order"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IsAdmin checks if the user is an admin
func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve IsAdmin value from the context
		isAdmin, exists := c.Get("IsAdmin")
		if !exists || (isAdmin != true && isAdmin != 1) { // Handle both boolean and integer types
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have admin privileges"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IsItMyProfile checks if the user is trying to access their own profile
func IsItMyProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assume user ID is set in the context after successful login
		userID := c.MustGet("id").(uint)

		// Extract the profile ID from the route parameter (assuming it's in the URL)
		profileID, err := strconv.ParseUint(c.Param("id"), 10, 32) // assuming the profile ID is in the URL like /profile/:id
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
			c.Abort()
			return
		}

		// Compare the IDs
		if uint(profileID) != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this profile"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ErrorHandler(c *gin.Context) {
	c.Next() // execute all the handlers

	// Check if there was an error
	if len(c.Errors) > 0 {
		for _, e := range c.Errors {
			log.Printf("Error: %v", e.Error())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
