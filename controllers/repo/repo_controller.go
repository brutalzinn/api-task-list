package repo_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	repo_service "github.com/brutalzinn/api-task-list/services/database/repo"
	authentication_util "github.com/brutalzinn/api-task-list/services/utils/authentication"
	hypermedia_util "github.com/brutalzinn/api-task-list/services/utils/hypermedia"
	"github.com/go-chi/chi/v5"
)

// @Summary      Get repo by id
// @Description  Get repo by id for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	repo, err := repo_service.Get(int64(id))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	links := map[string]any{}
	var repoMap = dto.ToRepoDTO(repo)
	hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/repo/%d", repo.ID), "DELETE")
	hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/repo/%d", repo.ID), "PUT")
	hypermedia_util.CreateHyperMedia(links, "detail", fmt.Sprintf("/repo/%d", repo.ID), "GET")
	repoMap.Links = links
	resp := response_entities.GenericResponse{
		Data: repoMap,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Update a repo
// @Description  Update a repo for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	user_id := authentication_util.GetCurrentUser(w, r)
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		log.Printf("wron url format %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var repo database_entities.Repo
	err = json.NewDecoder(r.Body).Decode(&repo)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := repo_service.Update(id, user_id, repo)
	if err != nil {
		log.Printf("error on update register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows > 1 {
		log.Printf("updates on  %d", rows)
	}
	resp := response_entities.GenericResponse{
		Message: "Repos updated",
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      List Repos
// @Description  List Repos for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo [get]
func List(w http.ResponseWriter, r *http.Request) {
	repos, err := repo_service.GetAll()
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var repoList = dto.ToRepoListDTO(repos)
	for i, repo := range repoList {
		links := map[string]any{}
		hypermedia_util.CreateHyperMedia(links, "task_list", fmt.Sprintf("/task/paginate?page=1&limit=10&repo_id=%d&order=DESC", repo.ID), "GET")
		hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/repo/%d", repo.ID), "DELETE")
		hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/repo/%d", repo.ID), "PUT")
		hypermedia_util.CreateHyperMedia(links, "detail", fmt.Sprintf("/repo/%d", repo.ID), "GET")
		repo.Links = links
		repoList[i] = repo
	}
	data := struct {
		Repos []dto.RepoDTO `json:"repos"`
	}{
		Repos: repoList,
	}
	resp := response_entities.GenericResponse{
		Data: data,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Paginate Repos
// @Description  Paginate Repos for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo/paginate [get]
func Paginate(w http.ResponseWriter, r *http.Request) {
	user_id := authentication_util.GetCurrentUser(w, r)
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 1
	}
	//order: 1 ASC -1 DESC
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "ASC"
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	offset := (page - 1) * limit
	repos, err := repo_service.Paginate(limit, offset, order, user_id)
	if err != nil {
		log.Printf("error on decode paginate json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totalTasks, _ := repo_service.Count()
	totalPages := (totalTasks + limit - 1) / limit
	currentPage := page

	var repoList = dto.ToRepoListDTO(repos)
	for i, repo := range repoList {
		links := map[string]any{}
		hypermedia_util.CreateHyperMedia(links, "task_list", fmt.Sprintf("/task/paginate?page=1&limit=10&repo_id=%d&order=DESC", repo.ID), "GET")
		hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/repo/%d", repo.ID), "DELETE")
		hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/repo/%d", repo.ID), "PUT")
		hypermedia_util.CreateHyperMedia(links, "detail", fmt.Sprintf("/repo/%d", repo.ID), "GET")
		repo.Links = links
		repoList[i] = repo
	}
	data := struct {
		Repos       []dto.RepoDTO `json:"repos"`
		TotalItems  int           `json:"totalItems"`
		TotalPages  int64         `json:"totalPages"`
		CurrentPage int64         `json:"currentPage"`
	}{
		Repos:       repoList,
		TotalItems:  len(repoList),
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}
	resp := response_entities.GenericResponse{
		Data: data,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Delete a repo
// @Description  Delete a repo for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	user_id := authentication_util.GetCurrentUser(w, r)
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := repo_service.Delete(id, user_id)
	if err != nil {
		log.Printf("error on delete register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Message: "Repo deleted",
	}
	if rows == 0 {
		resp = response_entities.GenericResponse{
			Message: "Cant delete this repo",
		}
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Create a repo
// @Description  Create a repo for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo [post]
func Create(w http.ResponseWriter, r *http.Request) {
	user_id := authentication_util.GetCurrentUser(w, r)
	var repo database_entities.Repo
	err := json.NewDecoder(r.Body).Decode(&repo)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	repo.UserId = user_id
	id, err := repo_service.Insert(repo)
	repo.ID = id
	resp := response_entities.GenericResponse{
		Message: fmt.Sprintf("Repo created %d", id),
		Data:    repo,
	}
	if err != nil {
		resp = response_entities.GenericResponse{
			Message: "Cant create this repo",
		}
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
