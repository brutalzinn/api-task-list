package dto

import (
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

type ApiKeyDTO struct {
	ID       string                      `json:"id"`
	Name     string                      `json:"name"`
	ExpireAt string                      `json:"expireAt"`
	Scopes   string                      `json:"scopes"`
	UpdateAt *string                     `json:"updateAt"`
	CreateAt *string                     `json:"createAt"`
	Links    []hypermedia.HypermediaLink `json:"links"`
}

func ToApiKeyDTO(apiKey database_entities.ApiKey) ApiKeyDTO {
	return ApiKeyDTO{
		ID:       apiKey.ID,
		Name:     apiKey.Name,
		ExpireAt: apiKey.ExpireAt,
		Scopes:   apiKey.Scopes,
		CreateAt: apiKey.CreateAt,
		UpdateAt: apiKey.UpdateAt,
	}
}

func ToApiKeyListDTO(original_apikeys []database_entities.ApiKey) []ApiKeyDTO {
	new_apikeys := make([]ApiKeyDTO, 0)
	for _, item := range original_apikeys {
		apiKeyDTO := ToApiKeyDTO(item)
		new_apikeys = append(new_apikeys, apiKeyDTO)
	}
	return new_apikeys
}
