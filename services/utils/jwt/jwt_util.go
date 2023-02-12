package jwt_util

import (
	"api-auto-assistant/configs"
	entities "api-auto-assistant/models"
	"errors"
	"fmt"
	"strings"
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
func VerifyJWT(tokenJWT string) (*jwt.Token, error) {
	var secretKey = configs.GetAuthSecret()
	safeSplit := strings.Split(tokenJWT, " ")
	fmt.Printf("token split count %d", len(safeSplit))
	if len(safeSplit) != 2 {
		return nil, errors.New("No valid token provided")
	}
	token := safeSplit[1]
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.New("No valid token provided")
	}
	fmt.Printf("claims %d", tkn.Claims)
	return tkn, err
}
