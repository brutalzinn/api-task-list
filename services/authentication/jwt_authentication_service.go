package authentication_service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/brutalzinn/api-task-list/common/scope"
	"github.com/brutalzinn/api-task-list/configs"
	request_entities "github.com/brutalzinn/api-task-list/models/request"
	user_service "github.com/brutalzinn/api-task-list/services/database/user"
	crypt_util "github.com/brutalzinn/api-task-list/utils/crypt"

	"github.com/golang-jwt/jwt/v4"
)

func Authentication(email string, password string) (userId string, err error) {
	user, err := user_service.FindByEmail(email)
	if err != nil {
		return "", errors.New("Invalid user")
	}
	validPassword := crypt_util.CheckPasswordHash(password, user.Password)
	if validPassword == false {
		return "", errors.New("Invalid user")
	}
	return user.ID, nil
}

func GenerateJWT(userId string) (string, error) {
	var secretKey = configs.GetAuthSecret()
	var authConfig = configs.GetConfig()
	expirationTime := time.Now().Add(time.Duration(authConfig.Authentication.Expiration) * time.Second).Unix()
	claims := request_entities.Claims{
		ID: userId,
		StandardClaims: jwt.StandardClaims{
			Subject:   userId,
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}
func VerifyJWT(tokenJWT string) (*request_entities.Claims, error) {
	var secretKey = configs.GetAuthSecret()
	safeSplit := strings.Split(tokenJWT, " ")
	///this messages will be removed after put this on mongoDB logger.
	if len(safeSplit) != 2 {
		return nil, errors.New("No valid token provided")
	}
	token := safeSplit[1]
	claims := request_entities.Claims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.New("No valid token provided")
	}
	return &claims, err
}
func GetCurrentUser(ctx context.Context) (user_id string, err error) {
	user_id, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", errors.New("No authorized to use this route")
	}
	return user_id, nil
}

func VerifyScope(ctx context.Context, requiredScopes []string) error {
	scopeClaim, ok := ctx.Value("scopes").(string)
	if !ok {
		return errors.New("No authorized to use this route")
	}
	scope := scope.New(scopeClaim)
	valid := scope.HasScope(requiredScopes)
	if !valid {
		return errors.New("No authorized to use this route")
	}
	return nil
}
