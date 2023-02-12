package task_controller

import (
	entities "api-auto-assistant/models"
	task_service "api-auto-assistant/services/database/task"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @Summary      Get task by id
// @Description  Get task by id for current user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} entities.Task
// @Router       /tasks/{id} [get]
func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	tasks, err := task_service.Get(int64(id))
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"Error":   false,
		"Message": tasks,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Update a task
// @Description  Update a task for current user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} entities.Task
// @Router       /tasks/{id} [put]
func Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("wron url format %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var task entities.Task
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
	if rows > 1 {
		log.Printf("updates on  %d", rows)
	}
	resp := map[string]any{
		"Error":   false,
		"Message": "Tasks updated",
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      List tasks
// @Description  List tasks for current user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200  {object} entities.Task
// @Router       /tasks [get]
func List(w http.ResponseWriter, r *http.Request) {
	tasks, err := task_service.GetAll()
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp := map[string]any{
		"Error":   false,
		"Message": tasks,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Delete a task
// @Description  Delete a task for current user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param id path int true "ID"
// @Success      200  {object} entities.Task
// @Router       /tasks/{id} [delete]
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

// @Summary      Create a task
// @Description  Create a task for current user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200  {object} entities.Task
// @Router       /tasks [post]
func Create(w http.ResponseWriter, r *http.Request) {
	var task entities.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("error on decode json %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id, err := task_service.Insert(task)
	var resp map[string]any
	if err != nil {
		resp = map[string]any{
			"Error":   true,
			"Message": fmt.Sprintf("Task creation failed %v", err),
		}
	} else {
		resp = map[string]any{
			"Error":   false,
			"Message": fmt.Sprintf("Task create %d", id),
		}
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
