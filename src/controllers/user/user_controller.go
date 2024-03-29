package user_controller

import (
	"encoding/json"
	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
	"log"
	"net/http"
	"strconv"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	request_entities "github.com/brutalzinn/api-task-list/models/request"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	authentication_util "github.com/brutalzinn/api-task-list/services/authentication"
	user_service "github.com/brutalzinn/api-task-list/services/database/user"
	crypt_util "github.com/brutalzinn/api-task-list/utils/crypt"
	"github.com/go-chi/chi/v5"
)

func Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	users, err := user_service.Get(id)
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
	id := chi.URLParam(r, "id")
	rows, err := user_service.Delete(id)
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
		"Message": "User deleted",
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
	hash, _ := crypt_util.HashPassword(user.Password, 15)
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
	var auth request_entities.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userId, err := authentication_service.Authentication(auth.Email, auth.Password)
	if err != nil {
		log.Printf("error on auth %v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	jwtToken, err := authentication_util.GenerateJWT(userId)
	resp := request_entities.AuthResponse{
		AccessToken: jwtToken,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
