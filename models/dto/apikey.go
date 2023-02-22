package dto

import (
	database_entities "github.com/brutalzinn/api-task-list/models/database"
	rest_entities "github.com/brutalzinn/api-task-list/models/rest"
)

type ApiKeyDTO struct {
	ID       int64                          `json:"id"`
	Name     string                         `json:"name"`
	ExpireAt string                         `json:"expireAt"`
	Scopes   string                         `json:"scopes"`
	UpdateAt *string                        `json:"updateAt"`
	CreateAt *string                        `json:"createAt"`
	Links    []rest_entities.HypermediaLink `json:"links"`
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

func ToApiKeyListDTO(original_apikeys []database_entities.ApiKey) (new_apikeys []ApiKeyDTO) {
	for _, item := range original_apikeys {
		apiKeyDTO := ToApiKeyDTO(item)
		new_apikeys = append(new_apikeys, apiKeyDTO)
	}
	return new_apikeys
}
