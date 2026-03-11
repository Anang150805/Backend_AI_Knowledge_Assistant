package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}