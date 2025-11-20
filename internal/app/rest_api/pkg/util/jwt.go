package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func SetupJWT(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateJWT(id int, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
