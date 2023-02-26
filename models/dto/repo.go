package dto

import (
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

type RepoDTO struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CreateAt    *string `json:"create_at"`
	UpdateAt    *string `json:"update_at"`
	Links       any     `json:"links"`
}

func ToRepoDTO(repo database_entities.Repo) RepoDTO {
	return RepoDTO{
		ID:          repo.ID,
		Title:       repo.Title,
		Description: repo.Description,
		CreateAt:    repo.CreateAt,
		UpdateAt:    repo.UpdateAt,
	}
}

func ToRepoListDTO(original_repos []database_entities.Repo) []RepoDTO {
	new_repos := make([]RepoDTO, 0)
	for _, item := range original_repos {
		repoDto := ToRepoDTO(item)
		new_repos = append(new_repos, repoDto)
	}
	return new_repos
}
