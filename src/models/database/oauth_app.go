package database_entities

import "time"

type OAuthApp struct {
	ID           int64
	AppName      string
	Mode         int64
	UserId       string
	ClientSecret string
	ClientId     string
	CreateAt     *time.Time
	UpdateAt     *time.Time
}
