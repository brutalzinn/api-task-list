package dto

import (
	"time"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

type TaskDTO struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Text        string     `json:"text"`
	Description string     `json:"description"`
	RepoId      int64      `json:"repo_id"`
	CreateAt    *time.Time `json:"create_at"`
	UpdateAt    *time.Time `json:"update_at"`
	Links       any        `json:"links"`
}

func ToTaskDTO(task database_entities.Task) TaskDTO {
	return TaskDTO{
		ID:          task.ID,
		Text:        task.Text,
		RepoId:      task.RepoId,
		Title:       task.Title,
		Description: task.Description,
		CreateAt:    task.CreateAt,
		UpdateAt:    task.UpdateAt,
	}
}

func ToTaskListDTO(original_repos []database_entities.Task) []TaskDTO {
	new_repos := make([]TaskDTO, 0)
	for _, item := range original_repos {
		repoDto := ToTaskDTO(item)
		new_repos = append(new_repos, repoDto)
	}
	return new_repos
}
