package entities

type Task struct {
	ID          int64   `json:id`
	Title       string  `json:title`
	Description string  `json:description`
	RepoId      int64   `json:repo_id`
	CreateAt    *string `json:create_at`
	UpdateAt    *string `json:update_at`
}
