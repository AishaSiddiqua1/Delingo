package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware is a Gin middleware that checks for a valid JWT token.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// Token should be in the format: Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			c.Abort()
			return
		}

		// Parse token and extract claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method, it should match the one used when creating the token
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key used to sign the token (this should be the same key used when signing)
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user ID from token claims (assuming you included it in the token during login)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"].(float64) // User ID is stored as a float64 in JWT claims
			c.Set("userID", uint(userID))         // Store user ID in the context
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims are invalid"})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}
