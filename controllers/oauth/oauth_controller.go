package oauth_controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	oauth_api_server "github.com/brutalzinn/api-task-list/oauth"
	"github.com/go-session/session"
)

func Token(w http.ResponseWriter, r *http.Request) {
	srv := oauth_api_server.GetOauthServer()
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	srv := oauth_api_server.GetOauthServer()
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	r.Form = form

	store.Delete("ReturnUri")
	store.Save()

	err = srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/oauth/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/auth.html")

}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", r.Form.Get("email"))
		store.Save()

		w.Header().Set("Location", "/oauth/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func Test(w http.ResponseWriter, r *http.Request) {
	srv := oauth_api_server.GetOauthServer()
	token, err := srv.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}
