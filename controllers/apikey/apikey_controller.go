package apikey_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	request_entities "github.com/brutalzinn/api-task-list/models/request"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	apikey_service "github.com/brutalzinn/api-task-list/services/database/apikey"
	apikey_util "github.com/brutalzinn/api-task-list/services/utils/apikey"
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
	if count >= 3 {
		log.Printf("error on count api key register %v", err)
		resp := map[string]any{
			"Error":   false,
			"Message": "You cant generate more api keys. Revoke one first.",
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	name := apiKeyRequest.Name
	nameNormalized := apikey_util.Normalize(name)
	apiKeysSameAppName, _ := apikey_service.CountByUserAndName(user_id, nameNormalized)
	if apiKeysSameAppName >= 1 {
		resp := map[string]any{
			"Error":   false,
			"Message": "You already have a api key with same app name registred. Revoke the old api key or generate a new api key with app name different.",
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	newApiKey, _ := apikey_util.CreateApiKey(user_id, nameNormalized)
	hashKey, _ := crypt_util.HashPassword(newApiKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	days := apiKeyRequest.ExpireAt.Days()
	expireAt := time.Now().AddDate(0, 0, days)
	scopes := []string{"task_insert", "task_delete", "task_read", "task_update"}
	apikey := database_entities.ApiKey{
		ApiKey:         hashKey,
		Scopes:         strings.Join(scopes, ","),
		Name:           apiKeyRequest.Name,
		NameNormalized: nameNormalized,
		UserId:         user_id,
		ExpireAt:       expireAt.Format(time.DateTime),
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
		Data:    map[string]any{"accesstoken": newApiKey},
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
func Revoke(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := apikey_service.Delete(int64(id))
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
		hypermedia_util.CreateHyperMedia(links, "revoke", fmt.Sprintf("/apikey/revoke/%d", apiKey.ID), "POST")
		apiKey.Links = links
		apiKeyList[i] = apiKey
	}

	resp := response_entities.GenericResponse{
		Data: apiKeyList,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
