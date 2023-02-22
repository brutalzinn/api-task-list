package database_entities

type Repo struct {
	ID          int64   `json:id`
	Title       string  `json:title`
	Description string  `json:description`
	UserId      int64   `json:user_id`
	CreateAt    *string `json:create_at`
	UpdateAt    *string `json:update_at`
	Tasks       *[]Task `json:tasks`
}
