package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"julia-conversation-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth middleware validates the JWT token and sets the userID and raw token in the context
func Auth(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if !cfg.ValidationEnabled {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in the format 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]
		c.Set("jwt", tokenString)

		if !cfg.ValidationEnabled {
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		sub, _ := claims.GetSubject()
		c.Set("UserID", sub)
		c.Next()
	}
}
