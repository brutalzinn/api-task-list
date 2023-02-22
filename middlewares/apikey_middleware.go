package middlewares

import (
	"context"
	"net/http"

	apikey_service "github.com/brutalzinn/api-task-list/services/database/apikey"
	apikey_util "github.com/brutalzinn/api-task-list/services/utils/apikey"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
)

func ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("x-api-key")
		if authHeader == "" {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		decrypt, err := crypt_util.Decrypt(authHeader)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		user_id, appName, err := apikey_util.GetApiKeyInfo(decrypt)
		count, err := apikey_service.CountByUserAndName(user_id, appName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if count == 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(r.Context(), "user_id", user_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
