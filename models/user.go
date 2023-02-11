package entities

type User struct {
	ID            int64   `json:id`
	Email         string  `json:email`
	Password      string  `json:password`
	Username      string  `json:username`
	FirebaseToken string  `json:firebaseToken`
	Create_at     *string `json:create_at`
	Update_at     *string `json:update_at`
}
