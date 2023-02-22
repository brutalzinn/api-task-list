package dto

import (
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

type TaskDTO struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	RepoId      int64   `json:"repo_id"`
	CreateAt    *string `json:"create_at"`
	UpdateAt    *string `json:"update_at"`
	Links       any     `json:"links"`
}

func ToTaskDTO(repo database_entities.Task) TaskDTO {
	return TaskDTO{
		ID:          repo.ID,
		RepoId:      repo.RepoId,
		Title:       repo.Title,
		Description: repo.Description,
		CreateAt:    repo.CreateAt,
		UpdateAt:    repo.UpdateAt,
	}
}

func ToTaskListDTO(original_repos []database_entities.Task) (new_repos []TaskDTO) {
	for _, item := range original_repos {
		repoDto := ToTaskDTO(item)
		new_repos = append(new_repos, repoDto)
	}
	return new_repos
}
