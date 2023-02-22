package task_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	task_service "github.com/brutalzinn/api-task-list/services/database/task"

	"github.com/go-chi/chi/v5"
)

// @Summary      Get task by id
// @Description  Get task by id for current user
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	tasks, err := task_service.Get(int64(id))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Data: tasks,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Update a task
// @Description  Update a task for current user
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("wron url format %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var task database_entities.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := task_service.Update(int64(id), task)
	if err != nil {
		log.Printf("error on update register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Message: "Task updated",
	}
	if rows == 0 {
		resp = response_entities.GenericResponse{
			Message: "Cant update task",
		}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      List tasks
// @Description  List tasks for current user
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task [get]
func List(w http.ResponseWriter, r *http.Request) {
	tasks, err := task_service.GetAll()
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Data: tasks,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Delete a task
// @Description  Delete a task for current user
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := task_service.Delete(int64(id))
	if err != nil {
		log.Printf("error on delete register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Message: "Task deleted",
	}

	if rows == 0 {
		resp = response_entities.GenericResponse{
			Message: "Cant delete this Task",
		}
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Create a task
// @Description  Create a task for current user
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task [post]
func Create(w http.ResponseWriter, r *http.Request) {
	var task database_entities.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id, err := task_service.Insert(task)
	task.ID = id
	resp := response_entities.GenericResponse{
		Message: fmt.Sprintf("Task created %d", id),
		Data:    task,
	}
	if err != nil {
		resp = response_entities.GenericResponse{
			Message: "Cant create this task",
		}
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Paginate Tasks
// @Description  Paginate Tasks for current user
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task/paginate [get]
func Paginate(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	//order: 1 ASC -1 DESC
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "ASC"
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	repo_id, err := strconv.Atoi(r.URL.Query().Get("repo_id"))
	if err != nil {
		repo_id = -1
	}
	offset := (page - 1) * limit

	repos, err := task_service.Paginate(repo_id, limit, offset, order)
	if err != nil {
		log.Printf("error on decode paginate json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totalTasks, _ := task_service.Count()
	totalPages := (totalTasks + limit - 1) / limit
	currentPage := page
	data := struct {
		Tasks       []database_entities.Task `json:"repos"`
		TotalItems  int                      `json:"totalItems"`
		TotalPages  int                      `json:"totalPages"`
		CurrentPage int                      `json:"currentPage"`
	}{
		Tasks:       repos,
		TotalItems:  len(repos),
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}

	resp := response_entities.GenericResponse{
		Data: data,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
