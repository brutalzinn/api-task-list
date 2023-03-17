package database_entities

import "time"

type OAuthApp struct {
	ID            int64
	AppName       string
	Mode          int64
	UserId        string
	OauthSecret   string
	OAuthClientId string
	CreateAt      *time.Time
	UpdateAt      *time.Time
}
