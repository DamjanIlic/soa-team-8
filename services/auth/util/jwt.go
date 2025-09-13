package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Tajni ključ za potpisivanje JWT tokena (u produkciji uzmi iz env var)
var jwtKey = []byte("supersecretkey")

// GenerateJWT vraća JWT token za korisnika
func GenerateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT proverava token i vraća claims
func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
