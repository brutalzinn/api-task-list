package database_entities

type ApiKey struct {
	ID             int64
	ApiKey         string
	Scopes         string
	Name           string
	NameNormalized string
	UserId         int64
	CreateAt       *string
	UpdateAt       *string
	ExpireAt       string
}


