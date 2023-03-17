package database_entities

import "time"

type Repo struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	UserId      string     `json:"user_id"`
	CreateAt    *time.Time `json:"create_at"`
	UpdateAt    *time.Time `json:"update_at"`
}
