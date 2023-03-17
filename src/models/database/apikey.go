package database_entities

import "time"

type ApiKey struct {
	ID             string
	ApiKey         string
	Scopes         string
	Name           string
	NameNormalized string
	UserId         string
	CreateAt       *time.Time
	UpdateAt       *time.Time
	ExpireAt       time.Time
}
