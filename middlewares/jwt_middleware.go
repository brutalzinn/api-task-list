package middlewares

import (
	jwt_util "api-auto-assistant/services/utils/jwt"
	"context"
	"fmt"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Auth middleware called")
		ctx := r.Context()
		authHeaderApiKey := r.Header.Get("x-api-key")
		if authHeaderApiKey != "" {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		authHeaderBaerer := r.Header.Get("Authorization")
		token, err := jwt_util.VerifyJWT(authHeaderBaerer)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(r.Context(), "user_id", token.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
