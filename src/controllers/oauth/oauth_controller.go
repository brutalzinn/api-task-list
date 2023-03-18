package oauth_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/brutalzinn/api-task-list/common"
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	request_entities "github.com/brutalzinn/api-task-list/models/request"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
	oauth_service "github.com/brutalzinn/api-task-list/services/database/oauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-session/session"
)

// @Summary      Get token oauth key
// @Description  Get token oauth for application
// @Tags         Oauth
// @Accept       json
// @Produce      json
// @Router       /oauth/generate [get]
func Token(w http.ResponseWriter, r *http.Request) {
	srv := authentication_service.GetOauthServer()
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary      Authorize oauth
// @Description  Authorize oauth for application
// @Tags         Oauth
// @Accept       json
// @Produce      json
// @Router       /oauth/authorize [get]
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

// @Summary      AuthHandler oauth
// @Description  AuthHandler oauth for application
// @Tags         Oauth
// @Accept       json
// @Produce      json
// @Router       /oauth/auth [get]
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

// @Summary      Test oauth
// @Description  Test oauth for application
// @Tags         Oauth
// @Accept       json
// @Produce      json
// @Router       /oauth/test [get]
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
		"scopes":     token.GetScope(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

// @Summary      Generate oauth
// @Description  Generate oauth for application
// @Tags         Oauth
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.OAuthResponse
// @Router       /oauth/generate [post]
func Generate(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication_service.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var request request_entities.OauthGenerateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
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
		AppName:  request.ApplicationName,
		Mode:     0,
		UserId:   userId,
		ClientId: clientId,
	}
	err = oauth_service.CreateOauthForUser(newApp)
	if err != nil {
		log.Printf("Cant create credentials. %v", err)
		response_entities.GenericMessageError(w, r, "Cant create your credentials.")
		return
	}
	data := response_entities.OAuthResponse{
		ClientId:     clientId,
		ClientSecret: secretId,
	}
	response_entities.GenericOK(w, r, data)
}

// @Summary      Regenerate oauth
// @Description  Regenerate oauth for application
// @Tags         Oauth
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.OAuthResponse
// @Router       /oauth/regenerate [post]
func Regenerate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userId, err := authentication_service.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if id == "" {
		log.Printf("error on decode json %v", id)
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	var request request_entities.OauthGenerateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		Domain: request.Callback,
	}
	clientStore.Update(newClient)
	data := response_entities.OAuthResponse{
		ClientId:     id,
		ClientSecret: secretId,
	}
	response_entities.GenericOK(w, r, data)
}

// @Summary      List oauth
// @Description  List oauth for application
// @Tags         Oauth
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /oauth/list [post]
func List(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication_service.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	oauthapps, err := oauth_service.List(userId)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var appsList = dto.ToOAuthListDTO(oauthapps)
	for i, item := range appsList {
		links := hypermedia.CreateHyperMediaLinksFor(item.ClientId, r.Context())
		appsList[i].Links = links
	}
	response_entities.GenericOK(w, r, appsList)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userId, err := authentication_service.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if id == "" {
		log.Printf("error on decode json %v", id)
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	clientStore := authentication_service.GetClientStore()
	client, err := clientStore.GetByID(r.Context(), id)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Cant delete your credentials.")
		return
	}

	err = oauth_service.DeleteOauthForUser(client, userId)
	if err != nil {
		log.Printf("Cant create credentials. %v", err)
		response_entities.GenericMessageError(w, r, "Cant delete your credentials.")
		return
	}

	newClient := &models.Client{
		ID: client.GetID(),
	}
	clientStore.Remove(newClient)
	response_entities.GenericOK(w, r, "Deleted with success.")
}
