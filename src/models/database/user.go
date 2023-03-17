package database_entities

import "time"

type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	Username      string     `json:"username"`
	FirebaseToken string     `json:"firebaseToken"`
	CreateAt      *time.Time `json:"create_at"`
	UpdateAt      *time.Time `json:"update_at"`
}
