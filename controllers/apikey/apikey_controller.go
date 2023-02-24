package apikey_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brutalzinn/api-task-list/configs"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	request_entities "github.com/brutalzinn/api-task-list/models/request"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	apikey_service "github.com/brutalzinn/api-task-list/services/database/apikey"
	apikey_util "github.com/brutalzinn/api-task-list/services/utils/apikey"
	converter_util "github.com/brutalzinn/api-task-list/services/utils/converter"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
	hypermedia_util "github.com/brutalzinn/api-task-list/services/utils/hypermedia"
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
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(int64)
	maxApiKeys := configs.GetApiConfig().MaxApiKeys
	if !ok {
		log.Printf("Error. usr dont authenticate and try to generate api key %d", user_id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var apiKeyRequest request_entities.ApiKeyRequest
	err := json.NewDecoder(r.Body).Decode(&apiKeyRequest)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	count, err := apikey_service.Count(user_id)
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
	nameNormalized := apikey_util.Normalize(name)
	apiKeysSameAppName, _ := apikey_service.CountByUserAndName(user_id, nameNormalized)
	if apiKeysSameAppName >= 1 {
		resp := response_entities.GenericResponse{
			Error:   true,
			Message: "You already have a api key with same app name registred. Revoke the old api key or generate a new api key with app name different.",
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	days := apiKeyRequest.ExpireAt.Days()
	expireAt := time.Now().AddDate(0, 0, days)
	expireAtFormat := expireAt.Format(time.RFC3339)
	uuid := apikey_util.CreateUUID()
	apiKeyCrypt, _ := apikey_util.CreateApiHash(user_id, nameNormalized, uuid, expireAtFormat)
	apiKeyHash, _ := crypt_util.HashPassword(apiKeyCrypt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	scopes := []string{"task_insert", "task_delete", "task_read", "task_update"}
	apikey := database_entities.ApiKey{
		ApiKey:         apiKeyHash,
		Scopes:         strings.Join(scopes, ","),
		Name:           apiKeyRequest.Name,
		NameNormalized: nameNormalized,
		UserId:         user_id,
		ExpireAt:       expireAtFormat,
	}

	_, err = apikey_service.Insert(apikey)
	if err != nil {
		log.Printf("error on update api key register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Error:   false,
		Message: "Api key generated",
		Data:    map[string]any{"api_key": apiKeyCrypt},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Generate api key
// @Description  Generate api key for user
// @Tags         ApiKeys
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /apikey/generate [post]
func Regenerate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(int64)
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !ok {
		log.Printf("Error. usr dont authenticate and try to regenerate api key %d", user_id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	apiKey, err := apikey_service.GetByIdAndUser(id, user_id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	date, err := converter_util.ToDateTime(apiKey.ExpireAt)
	expireAtFormat := date.Format(time.RFC3339)
	uuid := apikey_util.CreateUUID()
	nameNormalized := apiKey.NameNormalized
	apiKeyCrypt, _ := apikey_util.CreateApiHash(user_id, nameNormalized, uuid, expireAtFormat)
	apiKeyHash, _ := crypt_util.HashPassword(apiKeyCrypt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	apikey := database_entities.ApiKey{
		ApiKey:         apiKeyHash,
		Scopes:         apiKey.Scopes,
		Name:           apiKey.Name,
		NameNormalized: nameNormalized,
		UserId:         user_id,
		ExpireAt:       expireAtFormat,
	}

	_, err = apikey_service.Insert(apikey)
	if err != nil {
		log.Printf("error on update api key register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Error:   false,
		Message: "Api key regenerated.",
		Data:    map[string]any{"api_key": apiKeyCrypt},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(int64)
	if !ok {
		log.Printf("Error. usr dont authenticate and try to generate api key %d", user_id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := apikey_service.DeleteByIdAndUser(id, user_id)
	if err != nil {
		log.Printf("error on update register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	resp := response_entities.GenericResponse{}
	if rows == 0 {
		resp = response_entities.GenericResponse{
			Error:   true,
			Message: "Api key cant be revoked.",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp = response_entities.GenericResponse{
		Error:   false,
		Message: "Api key revoked.",
	}
	json.NewEncoder(w).Encode(resp)
}

// @Summary      List apikeys
// @Description  List apikeys for current user
// @Tags         ApiKeys
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /apikey [get]
func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(int64)
	if !ok {
		log.Printf("Error. usr dont authenticate and try to list api key %d", user_id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	apiKeys, err := apikey_service.GetAll(user_id)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var apiKeyList = dto.ToApiKeyListDTO(apiKeys)
	for i, apiKey := range apiKeyList {
		links := map[string]any{}
		//because i dont want to break web app here now.
		hypermedia_util.CreateHyperMedia(links, "regenerate", fmt.Sprintf("/apikey/regenerate/%d", apiKey.ID), "POST")
		hypermedia_util.CreateHyperMedia(links, "revoke", fmt.Sprintf("/apikey/delete/%d", apiKey.ID), "DELETE")
		hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/apikey/delete/%d", apiKey.ID), "DELETE")

		apiKey.Links = links
		apiKeyList[i] = apiKey
	}

	resp := response_entities.GenericResponse{
		Data: apiKeyList,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
