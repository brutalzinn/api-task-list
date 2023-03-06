package oauth_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/brutalzinn/api-task-list/common"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	request_entities "github.com/brutalzinn/api-task-list/models/request"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
	oauth_service "github.com/brutalzinn/api-task-list/services/database/oauth"
	"github.com/go-chi/chi"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-session/session"
)

func Token(w http.ResponseWriter, r *http.Request) {
	srv := authentication_service.GetOauthServer()
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	srv := authentication_service.GetOauthServer()
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
		common.OutputHTML(w, r, "static/auth.html")
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

	common.OutputHTML(w, r, "static/auth.html")

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

		userId, err := authentication_service.Authentication(email, password)
		if err != nil {
			common.OutputHTML(w, r, "static/error.html")
			return
		}
		store.Set("LoggedUserId", userId)
		store.Save()
		w.Header().Set("Location", "/oauth/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	common.OutputHTML(w, r, "static/login.html")
}

func Test(w http.ResponseWriter, r *http.Request) {
	srv := authentication_service.GetOauthServer()
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
	userId := authentication_service.GetCurrentUser(w, r)
	var request request_entities.OauthGenerateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	clientStore := authentication_service.GetClientStore()
	clientId := authentication_service.CreateUUID()
	secretId := authentication_service.CreateUUID()
	clientStore.Create(&models.Client{
		UserID: userId,
		ID:     clientId,
		Secret: secretId,
		Domain: request.Callback,
	})
	newApp := database_entities.OAuthApp{
		AppName:       request.ApplicationName,
		Mode:          0,
		UserId:        userId,
		OAuthClientId: clientId,
	}
	err = oauth_service.CreateOauthForUser(newApp)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Cant create your credentials.")
		return
	}
	data := response_entities.OAuthResponse{
		ClientId:     clientId,
		ClientSecret: secretId,
	}
	response_entities.GenericOK(w, r, data)
}
func RegenerateSecret(w http.ResponseWriter, r *http.Request) {
	userId := authentication_service.GetCurrentUser(w, r)
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	clientStore := authentication_service.GetClientStore()
	client, err := clientStore.GetByID(r.Context(), id)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Cant generate your credentials.")
		return
	}
	secretId := authentication_service.CreateUUID()
	newClient := &models.Client{
		UserID: userId,
		ID:     client.GetID(),
		Secret: secretId,
		Domain: client.GetDomain(),
	}
	clientStore.Update(newClient)
	data := response_entities.OAuthResponse{
		ClientId:     id,
		ClientSecret: secretId,
	}
	response_entities.GenericOK(w, r, data)
}
func Update(w http.ResponseWriter, r *http.Request) {
	userId := authentication_service.GetCurrentUser(w, r)
	var request request_entities.OauthGenerateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	clientStore := authentication_service.GetClientStore()
	client, err := clientStore.GetByID(r.Context(), id)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Cant update your credentials.")
		return
	}
	newClient := &models.Client{
		UserID: userId,
		ID:     client.GetID(),
		Secret: client.GetSecret(),
		Domain: request.Callback,
	}
	err = clientStore.Update(newClient)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Cant update your credentials.")
		return
	}
	newApp := database_entities.OAuthApp{
		AppName:       request.ApplicationName,
		Mode:          0,
		UserId:        userId,
		OAuthClientId: client.GetID(),
	}
	rows, err := oauth_service.Update(newApp)
	if rows == 0 {
		response_entities.GenericMessageError(w, r, "Cant create your credentials.")
		return
	}

	data := response_entities.OAuthResponse{
		ClientId: id,
	}
	response_entities.GenericOK(w, r, data)
}
func List(w http.ResponseWriter, r *http.Request) {
	userId := authentication_service.GetCurrentUser(w, r)
	oauthapps, err := oauth_service.List(userId)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response_entities.GenericOK(w, r, oauthapps)
}
