package jwt_util

import (
	"api-auto-assistant/configs"
	entities "api-auto-assistant/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(id int64) (string, error) {
	var secretKey = configs.GetAuthSecret()
	var expiration = configs.GetAuthExpiration()
	expirationTime := time.Now().Add(time.Duration(expiration) * time.Second)
	// Create the JWT claims, which includes the username and expiry time
	claims := entities.Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}
func VerifyJWT(token string) (*jwt.Token, error) {
	var secretKey = configs.GetAuthSecret()
	claims := entities.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	return tkn, err
}
