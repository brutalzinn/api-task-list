package repo_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	authentication_service "github.com/brutalzinn/api-task-list/services/authentication"
	repo_service "github.com/brutalzinn/api-task-list/services/database/repo"
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
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	repo, err := repo_service.Get(id)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var repoMap = dto.ToRepoDTO(repo)
	repoMap.Links = hypermedia.CreateHyperMediaLinksFor(repo.ID, r.Context())
	response_entities.GenericOK(w, r, repoMap)
}

// @Summary      Update a repo
// @Description  Update a repo for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo/{id} [put]
func Patch(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication_service.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
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
	rows, err := repo_service.Update(id, userId, repo)
	if err != nil {
		log.Printf("error on update register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		response_entities.GenericMessageError(w, r, "Cant update this repo")
		return
	}
	response_entities.GenericOK(w, r, "Repo updated")
}

// @Summary      Paginate Repos
// @Description  Paginate Repos for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo/paginate [get]
func Paginate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := authentication_service.GetCurrentUser(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// err = authentication_service.VerifyScope(ctx, []string{"read:repo", "create:repo"})
	// if err != nil {
	// 	log.Printf("error on read scopes %v", err)
	// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	// 	return
	// }
	currentPage, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		currentPage = 1
	}
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "ASC"
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	offset := (currentPage - 1) * limit
	repos, err := repo_service.Paginate(limit, offset, order, userId)
	if err != nil {
		log.Printf("error on decode paginate json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totalTasks, _ := repo_service.Count(userId)
	totalPages := (totalTasks + limit - 1) / limit
	var repoList = dto.ToRepoListDTO(repos)
	for i, repo := range repoList {
		links := hypermedia.CreateHyperMediaLinksFor(repo.ID, r.Context())
		repoList[i].Links = links
	}
	response_entities.PaginateRepo(w, r, repoList, totalPages, currentPage)
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
	userId, err := authentication_service.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := repo_service.Delete(id, userId)
	if err != nil {
		log.Printf("error on delete register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		response_entities.GenericMessageError(w, r, "Cant delete this repo")
		return
	}
	response_entities.GenericOK(w, r, "Repo deleted")
}

// @Summary      Create a repo
// @Description  Create a repo for current user
// @Tags         Repos
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /repo [post]
func Create(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication_service.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var repo database_entities.Repo
	err = json.NewDecoder(r.Body).Decode(&repo)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	repo.UserId = userId
	id, err := repo_service.Insert(repo)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Cant create this repo")
		return
	}
	response_entities.GenericOK(w, r, id)
}
