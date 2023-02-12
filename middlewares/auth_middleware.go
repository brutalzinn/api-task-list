package middlewares

import (
	"context"
	"fmt"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	fmt.Printf("accept middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user", "123")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
