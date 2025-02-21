package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(adminOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Token"})
		}

		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server misconfiguration"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		userID := uint(userIDFloat)

		if adminOnly {
			isAdmin, ok := claims["is_admin"].(bool)
			if !ok || !isAdmin {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
				return
			}
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
