package task_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	task_service "github.com/brutalzinn/api-task-list/services/database/task"
	hypermedia_util "github.com/brutalzinn/api-task-list/services/utils/hypermedia"

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
	task, err := task_service.Get(int64(id))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	links := map[string]any{}
	hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/task/%d", task.ID), "DELETE")
	hypermedia_util.CreateHyperMedia(links, "update_all", fmt.Sprintf("/task/%d", task.ID), "PUT")
	hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/task/%d", task.ID), "PATCH")
	taskDto := dto.ToTaskDTO(task)
	taskDto.Links = links
	resp := response_entities.GenericResponse{
		Data: taskDto,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Updates tasks
// @Description  Updates array of tasks
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task/{id} [put]
func Patch(w http.ResponseWriter, r *http.Request) {
	var tasks []database_entities.Task
	err := json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = task_service.UpdateTasks(tasks)
	if err != nil {
		log.Printf("error on update  tasks register %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Message: "All tasks updated",
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Replace all tasks
// @Description  Replace all tasks for a repo
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task/{id} [put]
func Put(w http.ResponseWriter, r *http.Request) {
	repo_id, err := strconv.ParseInt(r.URL.Query().Get("repo_id"), 10, 64)
	if err != nil {
		repo_id = -1
	}
	var tasks []database_entities.Task
	err = json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = task_service.ReplaceTasksByRepo(tasks, repo_id)
	if err != nil {
		log.Printf("error on replace data %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := response_entities.GenericResponse{
		Message: "All tasks replaced",
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
	var taskList = dto.ToTaskListDTO(tasks)
	for i, task := range taskList {
		links := make(map[string]any)
		hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/task/%d", task.ID), "DELETE")
		hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/task/%d", task.ID), "PATCH")
		hypermedia_util.CreateHyperMedia(links, "update_all", fmt.Sprintf("/task/%d", task.ID), "PUT")
		hypermedia_util.CreateHyperMedia(links, "detail", fmt.Sprintf("/task/%d", task.ID), "GET")
		task.Links = links
		taskList[i] = task
	}
	data := struct {
		Tasks []dto.TaskDTO `json:"tasks"`
	}{
		Tasks: taskList,
	}
	resp := response_entities.GenericResponse{
		Data: data,
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

	tasks, err := task_service.Paginate(repo_id, limit, offset, order)
	if err != nil {
		log.Printf("error on decode paginate json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totalTasks, _ := task_service.Count()
	totalPages := (totalTasks + limit - 1) / limit
	currentPage := page

	var taskList = dto.ToTaskListDTO(tasks)
	for i, task := range taskList {
		links := map[string]any{}
		hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/task/%d", task.ID), "DELETE")
		hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/task/%d", task.ID), "PATCH")
		hypermedia_util.CreateHyperMedia(links, "update_all", fmt.Sprintf("/task/%d", task.ID), "PUT")
		hypermedia_util.CreateHyperMedia(links, "detail", fmt.Sprintf("/task/%d", task.ID), "GET")
		task.Links = links
		taskList[i] = task
	}

	data := struct {
		Tasks       []dto.TaskDTO `json:"tasks"`
		TotalItems  int           `json:"totalItems"`
		TotalPages  int           `json:"totalPages"`
		CurrentPage int           `json:"currentPage"`
	}{
		Tasks:       taskList,
		TotalItems:  len(taskList),
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}

	resp := response_entities.GenericResponse{
		Data: data,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
