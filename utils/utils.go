// utils/utils.go
package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JSONErrorResponse sends a standardized JSON error response
func JSONErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// GenerateJWT generates a JWT token for a given user ID
func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiry time (24 hours)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key")) // Replace with your secret key
}

func GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 168).Unix(), // Refresh token expires in 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key")) // Replace with your secret key
}
