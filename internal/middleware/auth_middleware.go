package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/wahyuoi/sbc/internal/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		parser := jwt.Parser{
			UseJSONNumber: true,
		}
		token, err := parser.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			return config.GetJWTSecret(), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId, err := claims["user_id"].(json.Number).Int64()
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}
			c.Set("user_email", claims["email"])
			c.Set("user_id", userId)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
