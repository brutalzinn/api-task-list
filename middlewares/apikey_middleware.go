package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	apikey_service "github.com/brutalzinn/api-task-list/services/database/apikey"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
)

func ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Api Key middleware called")
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
		apikeyformat := strings.Split(decrypt, "-")
		if len(apikeyformat) != 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		user_id, _ := strconv.ParseInt(apikeyformat[0], 10, 64)
		apikeycrypt := apikeyformat[1]
		apikeyfound, err := apikey_service.Get(user_id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		apiKeyValid := crypt_util.CheckPasswordHash(apikeycrypt, apikeyfound.ApiKey)
		if apiKeyValid == false {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(r.Context(), "user_id", apikeyfound.UserId)
		fmt.Printf("Api Key middleware OK %s", decrypt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
