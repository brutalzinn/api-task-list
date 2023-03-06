package dto

import (
	"time"

	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

type OAuthDTO struct {
	ID            int64
	AppName       string
	Mode          int64
	UserId        string
	OAuthClientId string
	CreateAt      *time.Time
	UpdateAt      *time.Time
	Links         []hypermedia.HypermediaLink `json:"links"`
}

func ToOAuthDTO(oauth database_entities.OAuthApp) OAuthDTO {
	return OAuthDTO{
		ID:            oauth.ID,
		AppName:       oauth.AppName,
		Mode:          oauth.Mode,
		UserId:        oauth.UserId,
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
