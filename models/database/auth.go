package database_entities

import "github.com/golang-jwt/jwt/v4"

type AuthRequest struct {
	Email    string `json:email`
	Password string `json:password`
}
type AuthResponse struct {
	AccessToken string `json:access_token`
}
type Claims struct {
	ID int64 `json:id`
	jwt.RegisteredClaims
}
