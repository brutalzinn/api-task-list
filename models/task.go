package entities

type Task struct {
	ID          int64   `json:id`
	Title       string  `json:title`
	Description string  `json:description`
	Create_at   *string `json:create_at`
	Update_at   *string `json:update_at`
}
