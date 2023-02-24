package middlewares

import (
	"context"
	"net/http"

	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeaderApiKey := r.Header.Get("x-api-key")
		if authHeaderApiKey != "" {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		authHeaderBaerer := r.Header.Get("Authorization")
		token, err := authentication_service.VerifyJWT(authHeaderBaerer)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(r.Context(), "user_id", token.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
