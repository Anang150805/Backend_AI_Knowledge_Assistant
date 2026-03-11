package middleware

import (
	"net/http"
	"os"
	"strings"

	"backend-AI-Knowledge-Assistant/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "invalid claims")
			c.Abort()
			return
		}

		userID, ok1 := claims["user_id"].(string)
		username, ok2 := claims["username"].(string)
		if !ok1 || !ok2 {
			response.Error(c, http.StatusUnauthorized, "invalid token payload")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Next()
	}
}