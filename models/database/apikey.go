package database_entities

type ApiKey struct {
	ID             string
	ApiKey         string
	Scopes         string
	Name           string
	NameNormalized string
	UserId         string
	CreateAt       *string
	UpdateAt       *string
	ExpireAt       string
}
