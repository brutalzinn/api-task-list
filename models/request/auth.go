package request_entities

import "github.com/golang-jwt/jwt/v4"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	AccessToken string `json:"accesstoken"`
}
type Claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
