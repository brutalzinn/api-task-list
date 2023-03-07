package middlewares

import (
	"context"
	"net/http"

	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
)

func ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("x-api-key")
		if authHeader == "" {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		apiKey, err := authentication_service.VerifyApiKey(authHeader)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, "user_id", apiKey.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
