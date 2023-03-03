package database_entities

import "time"

type OAuthApp struct {
	AppName       string     `json:"appname"`
	Mode          int64      `json:"mode"`
	UserId        string     `json:"user_id"`
	OAuthClientId string     `json:"oauth_client_id"`
	CreateAt      *time.Time `json:"create_at"`
	UpdateAt      *time.Time `json:"update_at"`
}
