package jwt

import (
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(userID string) (string, error) {
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(JwtSecret)
}

func ValidateToken(token string) (*jwtlib.Token, map[string]interface{}, error) {
	claims := jwtlib.MapClaims{}
	t, err := jwtlib.ParseWithClaims(token, claims, func(t *jwtlib.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	return t, claims, err
}
