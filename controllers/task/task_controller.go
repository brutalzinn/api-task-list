package task_controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	"github.com/brutalzinn/api-task-list/models/dto"
	response_entities "github.com/brutalzinn/api-task-list/models/response"
	repo_service "github.com/brutalzinn/api-task-list/services/database/repo"
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
	task, err := task_service.Get(int64(id))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	links := map[string]any{}
	// hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/task/%d", task.ID), "DELETE")
	// hypermedia_util.CreateHyperMedia(links, "update_all", fmt.Sprintf("/task/%d", task.ID), "PUT")
	// hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/task/%d", task.ID), "PATCH")
	taskDto := dto.ToTaskDTO(task)
	taskDto.Links = links
	response_entities.GenericOK(w, r, taskDto)
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
		response_entities.GenericOK(w, r, "Cant update tasks")
		return
	}
	response_entities.GenericOK(w, r, "All tasks updated")
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
	repoId, err := strconv.ParseInt(r.URL.Query().Get("repo_id"), 10, 64)
	if err != nil {
		repoId = -1
	}
	var tasks []database_entities.Task
	err = json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = repo_service.Get(repoId)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Repo not found")
		return
	}
	err = task_service.ReplaceTasksByRepo(tasks, repoId)
	if err != nil {
		log.Printf("error on replace data %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response_entities.GenericOK(w, r, "All tasks replaced")
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
	if rows == 0 {
		response_entities.GenericMessageError(w, r, "Cant delete this Task")
		return
	}
	response_entities.GenericOK(w, r, "Task deleted")
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
	_, err = repo_service.Get(task.RepoId)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Repo dont found")
		return
	}
	id, err := task_service.Insert(task)
	if err != nil {
		response_entities.GenericMessageError(w, r, "Cant create this task")
		return
	}
	response_entities.GenericOK(w, r, id)
}

// @Summary      Paginate Tasks
// @Description  Paginate Tasks for current user
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Success      200  {object} response_entities.GenericResponse
// @Router       /task/paginate [get]
func Paginate(w http.ResponseWriter, r *http.Request) {
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
	repoId, err := strconv.ParseInt(r.URL.Query().Get("repo_id"), 10, 64)
	if err != nil {
		repoId = -1
	}
	offset := (currentPage - 1) * limit
	tasks, err := task_service.Paginate(repoId, limit, offset, order)
	if err != nil {
		log.Printf("error on decode paginate json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totalTasks, _ := task_service.Count(repoId)
	totalPages := (totalTasks + limit - 1) / limit
	taskList := dto.ToTaskListDTO(tasks)
	// for i, task := range taskList {
	// 	links := map[string]any{}
	// 	hypermedia_util.CreateHyperMedia(links, "delete", fmt.Sprintf("/task/%d", task.ID), "DELETE")
	// 	hypermedia_util.CreateHyperMedia(links, "update_one", fmt.Sprintf("/task/%d", task.ID), "PATCH")
	// 	hypermedia_util.CreateHyperMedia(links, "update_all", fmt.Sprintf("/task?repo_id=%d", task.RepoId), "PUT")
	// 	hypermedia_util.CreateHyperMedia(links, "detail", fmt.Sprintf("/task/%d", task.ID), "GET")
	// 	task.Links = links
	// 	taskList[i] = task
	// }
	response_entities.PaginateTask(w, r, taskList, totalPages, currentPage)
}
