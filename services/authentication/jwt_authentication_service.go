package authentication_service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/brutalzinn/api-task-list/configs"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	user_service "github.com/brutalzinn/api-task-list/services/database/user"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"

	"github.com/golang-jwt/jwt/v4"
)

func Authentication(email string, password string) (user database_entities.User, err error) {
	user, err = user_service.FindByEmail(email)
	if err != nil {
		return user, errors.New("Invalid user")
	}
	validPassword := crypt_util.CheckPasswordHash(password, user.Password)
	if validPassword == false {
		return user, errors.New("Invalid user")
	}
	return user, nil
}

func GenerateJWT(id string) (string, error) {
	var secretKey = configs.GetAuthSecret()
	var authConfig = configs.GetConfig()
	expirationTime := time.Now().Add(time.Duration(authConfig.Authentication.Expiration) * time.Second)
	claims := database_entities.Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}
func VerifyJWT(tokenJWT string) (*database_entities.Claims, error) {
	var secretKey = configs.GetAuthSecret()
	safeSplit := strings.Split(tokenJWT, " ")
	if len(safeSplit) != 2 {
		return nil, errors.New("No valid token provided")
	}
	token := safeSplit[1]
	claims := database_entities.Claims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.New("No valid token provided")
	}
	return &claims, err
}
func GetCurrentUser(w http.ResponseWriter, r *http.Request) (user_id string) {
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	return
}
