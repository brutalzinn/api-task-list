package dto

import (
	"time"

	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

type OAuthDTO struct {
	ID            int64                       `json:"id"`
	AppName       string                      `json:"app_name"`
	Mode          int64                       `json:"mode"`
	OAuthClientId string                      `json:"client_id"`
	CreateAt      *time.Time                  `json:"create_at"`
	UpdateAt      *time.Time                  `json:"update_at"`
	Links         []hypermedia.HypermediaLink `json:"links"`
}

func ToOAuthDTO(oauth database_entities.OAuthApp) OAuthDTO {
	return OAuthDTO{
		ID:            oauth.ID,
		AppName:       oauth.AppName,
		Mode:          oauth.Mode,
		OAuthClientId: oauth.OAuthClientId,
		CreateAt:      oauth.CreateAt,
		UpdateAt:      oauth.UpdateAt,
	}
}

func ToOAuthListDTO(original_repos []database_entities.OAuthApp) []OAuthDTO {
	new_repos := make([]OAuthDTO, 0)
	for _, item := range original_repos {
		authDTO := ToOAuthDTO(item)
		new_repos = append(new_repos, authDTO)
	}
	return new_repos
}
