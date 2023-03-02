package apikey_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/brutalzinn/api-task-list/configs"
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	request_entities "github.com/brutalzinn/api-task-list/models/request"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
	authentication_util "github.com/brutalzinn/api-task-list/services/authentication"
	apikey_service "github.com/brutalzinn/api-task-list/services/database/apikey"
	converter_util "github.com/brutalzinn/api-task-list/services/utils/converter"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
	"github.com/go-chi/chi/v5"
)

// @Summary      Generate api key
// @Description  Generate api key for user
// @Tags         ApiKeys
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /apikey/generate [post]
func Generate(w http.ResponseWriter, r *http.Request) {
	userId := authentication_util.GetCurrentUser(w, r)
	maxApiKeys := configs.GetConfig().API.MaxApiKeys
	var apiKeyRequest request_entities.ApiKeyRequest
	err := json.NewDecoder(r.Body).Decode(&apiKeyRequest)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	count, err := apikey_service.Count(userId)
	if err != nil {
		log.Printf("error on count api key register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if count >= maxApiKeys {
		log.Printf("error on count api key register %v", err)
		resp := response_entities.GenericResponse{
			Error:   true,
			Message: fmt.Sprintf("You cant generate more than %d api keys. Revoke one first.", maxApiKeys),
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	name := apiKeyRequest.Name
	nameNormalized := converter_util.Normalize(name)
	apiKeysSameAppName, _ := apikey_service.CountByUserAndName(userId, nameNormalized)
	if apiKeysSameAppName >= 1 {
		response_entities.GenericMessageError(w, r, "You already have a api key with same app name registred. Revoke the old api key or generate a new api key with app name different.")
		return
	}
	days := apiKeyRequest.ExpireAt.Days()
	expireAt := time.Now().Add(time.Hour * 24 * time.Duration(days))
	uuid := authentication_service.CreateUUID()
	randomFactor := authentication_service.CreateRandomFactor()
	apiKeyCrypt, _ := authentication_service.CreateApiKeyCrypt(uuid, randomFactor)
	newApiKey := authentication_service.CreateApiPrefix(apiKeyCrypt, nameNormalized)
	apiKeyHash, err := crypt_util.HashPassword(newApiKey, 4)
	if err != nil {
		///temporary because we have a limit of 72 bytes. We need change the hash algorithm or remove the nameNormalized for hash generation3
		response_entities.GenericMessageError(w, r, "Name needs be max 11 length")
		return
	}
	scopes := []string{"task_manager", "repo_manager"}
	apikey := database_entities.ApiKey{
		ID:             uuid,
		ApiKey:         apiKeyHash,
		Scopes:         strings.Join(scopes, ","),
		Name:           apiKeyRequest.Name,
		NameNormalized: nameNormalized,
		UserId:         userId,
		ExpireAt:       expireAt,
	}
	_, err = apikey_service.Insert(apikey)
	if err != nil {
		log.Printf("error on update api key register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response_entities.GenericOK(w, r, map[string]any{"api_key": newApiKey})
}

// @Summary      Generate api key
// @Description  Generate api key for user
// @Tags         ApiKeys
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /apikey/generate [post]
func Regenerate(w http.ResponseWriter, r *http.Request) {
	userId := authentication_util.GetCurrentUser(w, r)
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	apiKey, err := apikey_service.GetByIdAndUser(id, userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	var expireDiff time.Duration
	if apiKey.UpdateAt != nil {
		expireDiff = apiKey.ExpireAt.Sub(*apiKey.UpdateAt)
	} else {
		expireDiff = apiKey.ExpireAt.Sub(*apiKey.CreateAt)
	}
	apiKey.ExpireAt = time.Now().Add(expireDiff)
	nameNormalized := apiKey.NameNormalized
	randomFactor := authentication_service.CreateRandomFactor()
	apiKeyCrypt, _ := authentication_service.CreateApiKeyCrypt(apiKey.ID, randomFactor)
	newApiKey := authentication_service.CreateApiPrefix(apiKeyCrypt, nameNormalized)
	apiKeyHash, _ := crypt_util.HashPassword(newApiKey, 4)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	apikey := database_entities.ApiKey{
		ApiKey:         apiKeyHash,
		Scopes:         apiKey.Scopes,
		Name:           apiKey.Name,
		NameNormalized: nameNormalized,
		UserId:         userId,
		ExpireAt:       apiKey.ExpireAt,
	}
	_, err = apikey_service.Update(id, apikey)
	if err != nil {
		log.Printf("error on update api key register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response_entities.GenericOK(w, r, map[string]any{"api_key": newApiKey})
}

// @Summary      Revoke apikey
// @Description  Revoke a user apikey
// @Tags         ApiKeys
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /apikey/revoke/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	userId := authentication_util.GetCurrentUser(w, r)
	id := chi.URLParam(r, "id")
	rows, err := apikey_service.DeleteByIdAndUser(id, userId)
	if err != nil {
		log.Printf("error on update register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		response_entities.GenericMessageError(w, r, "Something happened. Api key cant be revoked.")
		return
	}
	response_entities.GenericOK(w, r, "Api key revoked.")
}

// @Summary      List apikeys
// @Description  List apikeys for current user
// @Tags         ApiKeys
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /apikey [get]
func List(w http.ResponseWriter, r *http.Request) {
	userId := authentication_util.GetCurrentUser(w, r)
	apiKeys, err := apikey_service.GetAll(userId)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var apiKeyList = dto.ToApiKeyListDTO(apiKeys)
	ctx := r.Context()
	links, _ := ctx.Value("links").([]hypermedia.HypermediaLink)
	for i, repo := range apiKeyList {
		var hypermediaLink []hypermedia.HypermediaLink
		for _, link := range links {
			link.Href = fmt.Sprintf(link.Href, repo.ID)
			hypermediaLink = append(hypermediaLink, link)
		}
		apiKeyList[i].Links = hypermediaLink
	}
	response_entities.GenericOK(w, r, apiKeyList)
}
