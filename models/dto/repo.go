package dto

import (
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	rest_entities "github.com/brutalzinn/api-task-list/models/rest"
)

type RepoDTO struct {
	ID          int64
	Title       string
	Description string
	UserId      int64
	CreateAt    *string
	UpdateAt    *string
	Links       []rest_entities.HypermediaLink
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

func ToRepoListDTO(original_repos []database_entities.Repo) (new_repos []RepoDTO) {
	for _, item := range original_repos {
		repoDto := ToRepoDTO(item)
		new_repos = append(new_repos, repoDto)
	}
	return new_repos
}
