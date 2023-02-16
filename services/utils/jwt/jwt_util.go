package jwt_util

import (
	"errors"
	"strings"
	"time"

	"github.com/brutalzinn/api-task-list/configs"
	entities "github.com/brutalzinn/api-task-list/models"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(id int64) (string, error) {
	var secretKey = configs.GetAuthSecret()
	var authConfig = configs.GetAuthConfig()
	expirationTime := time.Now().Add(time.Duration(authConfig.Expiration) * time.Second)
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
func VerifyJWT(tokenJWT string) (*entities.Claims, error) {
	var secretKey = configs.GetAuthSecret()
	safeSplit := strings.Split(tokenJWT, " ")
	if len(safeSplit) != 2 {
		return nil, errors.New("No valid token provided")
	}
	token := safeSplit[1]
	claims := entities.Claims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.New("No valid token provided")
	}
	return &claims, err
}
