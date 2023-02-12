package middlewares

import (
	jwt_util "api-auto-assistant/services/utils/jwt"
	"context"
	"fmt"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Auth middleware called")
		authHeader := r.Header.Get("Authorization")
		token, err := jwt_util.VerifyJWT(authHeader)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "id", "1")
		fmt.Printf("Auth middleware OK %s", token.Header)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
