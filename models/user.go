package entities

type User struct {
	ID            int64   `json:id`
	Email         string  `json:email`
	Password      string  `json:password`
	Username      string  `json:username`
	FirebaseToken string  `json:firebaseToken`
	CreateAt      *string `json:create_at`
	UpdateAt      *string `json:update_at`
}
