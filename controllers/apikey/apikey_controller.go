package apikey_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	entities "github.com/brutalzinn/api-task-list/models"
	apikey_service "github.com/brutalzinn/api-task-list/services/database/apikey"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"

	"github.com/google/uuid"
)

// @Summary      Generate api key
// @Description  Generate api key for user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} entities.Task
// @Router       /tasks/{id} [get]
func Generate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(int64)
	if !ok {
		log.Printf("Error. usr dont authenticate and try to generate api key %d", user_id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	scopes := []string{"task_insert", "task_delete", "task_read", "task_update"}
	count, err := apikey_service.Count(user_id)
	if err != nil {
		log.Printf("error on update api key register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		log.Printf("error on update api key register %v", err)
		resp := map[string]any{
			"Error":   false,
			"Message": "You cant generate a api key. Revoke your own api key first.",
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	uuid := uuid.New().String()
	uuidnormalized := strings.Replace(uuid, "-", "", -1)
	cryptKey, _ := crypt_util.HashPassword(uuidnormalized)
	keyhash, err := crypt_util.Encrypt(fmt.Sprintf("%d-%s", user_id, uuidnormalized))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	apikey := entities.ApiKey{
		ApiKey: cryptKey,
		Scopes: strings.Join(scopes, ","),
		UserId: user_id,
	}
	_, err = apikey_service.Insert(apikey)
	if err != nil {
		log.Printf("error on update api key register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"Error":       false,
		"Message":     "Api key generated",
		"AccessToken": keyhash,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Revoke apikey
// @Description  Revoke a user apikey
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} entities.Task
// @Router       /tasks/{id} [put]
func Revoke(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_id, ok := ctx.Value("user_id").(int64)
	if !ok {
		log.Printf("Error. usr dont authenticate and try to revoke api key %d", user_id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := apikey_service.Delete(int64(user_id))
	if err != nil {
		log.Printf("error on update register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	resp := map[string]any{}
	if rows == 0 {
		resp = map[string]any{
			"Error":   false,
			"Message": "Cant revoke api key.",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp = map[string]any{
		"Error":   false,
		"Message": "Api key revoked.",
	}
	json.NewEncoder(w).Encode(resp)
}
