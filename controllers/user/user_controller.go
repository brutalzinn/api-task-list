package user_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	user_service "github.com/brutalzinn/api-task-list/services/database/user"
	crypt_utils "github.com/brutalzinn/api-task-list/services/utils/crypt"
	jwt_util "github.com/brutalzinn/api-task-list/services/utils/jwt"

	"github.com/go-chi/chi/v5"
)

func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	users, err := user_service.Get(int64(id))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"Error":   false,
		"Message": users,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("wron url format %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var user database_entities.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := user_service.Update(int64(id), user)
	if err != nil {
		log.Printf("error on update register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows > 1 {
		log.Printf("updates on  %d", rows)
	}
	resp := map[string]any{
		"Error":   false,
		"Message": "users updated",
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func List(w http.ResponseWriter, r *http.Request) {
	users, err := user_service.GetAll()
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"Error":   false,
		"Message": users,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := user_service.Delete(int64(id))
	if err != nil {
		log.Printf("error on delete register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows > 1 {
		log.Printf("delete on  %d", rows)
	}
	resp := map[string]any{
		"Error":   false,
		"Message": "Task deleted",
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func Create(w http.ResponseWriter, r *http.Request) {
	var user database_entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	hash, _ := crypt_utils.HashPassword(user.Password)
	user.Password = hash
	_, err = user_service.Insert(user)
	resp := response_entities.GenericResponse{
		Error:   true,
		Message: "Registred with sucess",
	}
	if err != nil {
		resp = response_entities.GenericResponse{
			Error:   true,
			Message: "Cant register now",
		}
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func Login(w http.ResponseWriter, r *http.Request) {
	var auth database_entities.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	user, err := user_service.FindByEmail(auth.Email)
	validPassword := crypt_utils.CheckPasswordHash(auth.Password, user.Password)
	if validPassword == false {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	jwtToken, err := jwt_util.GenerateJWT(user.ID)
	resp := database_entities.AuthResponse{
		AccessToken: jwtToken,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
