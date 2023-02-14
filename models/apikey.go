package entities

type ApiKey struct {
	ID        int64   `json:id`
	ApiKey    string  `json:apikey`
	Scopes    string  `json:scopes`
	UserId    int64   `json:user_id`
	Create_at *string `json:create_at`
	Update_at *string `json:update_at`
}
