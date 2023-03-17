package response_entities

import (
	"net/http"

	"github.com/brutalzinn/api-task-list/models/dto"
)

func PaginateTask(w http.ResponseWriter, r *http.Request, list []dto.TaskDTO, totalPages int64, currentPage int64) {
	data := struct {
		List        []dto.TaskDTO `json:"tasks"`
		TotalItems  int           `json:"totalItems"`
		TotalPages  int64         `json:"totalPages"`
		CurrentPage int64         `json:"currentPage"`
	}{
		List:        list,
		TotalItems:  len(list),
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}
	GenericOK(w, r, data)
}
func PaginateRepo(w http.ResponseWriter, r *http.Request, list []dto.RepoDTO, totalPages int64, currentPage int64) {
	data := struct {
		List        []dto.RepoDTO `json:"repos"`
		TotalItems  int           `json:"totalItems"`
		TotalPages  int64         `json:"totalPages"`
		CurrentPage int64         `json:"currentPage"`
	}{
		List:        list,
		TotalItems:  len(list),
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}
	GenericOK(w, r, data)
}
