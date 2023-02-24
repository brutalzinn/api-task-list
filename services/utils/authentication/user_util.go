package authentication_util

import (
	"net/http"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) (user_id int64) {
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(int64)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	return
}
