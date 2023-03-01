package oauth_controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	request_entities "github.com/brutalzinn/api-task-list/models/request"
	oauth_api_server "github.com/brutalzinn/api-task-list/oauth"
	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
	"github.com/go-oauth2/oauth2/v4"
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
		w.WriteHeader(http.StatusUnauthorized)
		outputHTML(w, r, "static/auth.html")
	}
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedUserId"); !ok {
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
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		user, err := authentication_service.Authentication(email, password)
		if err != nil {
			outputHTML(w, r, "static/error.html")
			return
		}
		store.Set("LoggedUserId", user.ID)
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

func Generate(w http.ResponseWriter, r *http.Request) {
	var request request_entities.OauthGenerateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// userId := authentication_service.GetCurrentUser(w, r)
	// srv := oauth_api_server.GetOauthServer()
	tokenManager := oauth_api_server.GetTokenStore()
	clientId := authentication_service.CreateRandomFactor()
	secretId := authentication_service.CreateUUID()
	tokenInfo := oauth2.TokenInfo.New(nil)
	tokenManager.Create(context.TODO(), tokenInfo)

	// clientManager.Create(clientId, &models.Client{
	// 	ID:     clientId,
	// 	UserID: userId,
	// 	Secret: secretId,
	// 	Domain: request.Callback,
	// })

	data := map[string]interface{}{
		// "user_id":   token.GetUserID(),
		"client_id": clientId,
		"secret_id": secretId,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
