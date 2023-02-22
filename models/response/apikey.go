package response_entities

type ApiKeyResponse struct {
	Name     string  `json:"name"`
	ExpireAt string  `json:"expireAt"`
	CreateAt *string `json:"createAt"`
	Scopes   string  `json:"scopes"`
	UpdateAt *string `json:"updateAt"`
}
