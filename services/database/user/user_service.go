package user_service

import (
	"time"

	"github.com/brutalzinn/api-task-list/db"
	entities "github.com/brutalzinn/api-task-list/models"
)

func Delete(id int64) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func GetAll() (users []entities.User, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT * FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		var user entities.User
		err = rows.Scan(&user.ID, &user.Email, &user.Password,
			&user.Username, &user.FirebaseToken, &user.CreateAt,
			&user.UpdateAt)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return
}
func Get(id int64) (user entities.User, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT * FROM users WHERE id=$1", id)
	err = row.Scan(&user.ID, &user.Email, &user.Password,
		&user.Username, &user.FirebaseToken, &user.CreateAt,
		&user.UpdateAt)
	return
}

func Insert(user entities.User) (id int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	sql := "INSERT INTO users (email, password, username, firebaseToken, create_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = conn.QueryRow(sql, user.Email, user.Password, user.Username, user.FirebaseToken, time.Now()).Scan(&id)
	return
}
func Update(id int64, user entities.User) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("UPDATE users SET email=$1, password=$2, username=$3, firebaseToken=$4, update_at=$5 WHERE id=$6",
		user.Email, user.Password, user.Username, user.FirebaseToken, user.UpdateAt, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func FindByEmail(email string) (user entities.User, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT id, password FROM users WHERE email=$1", email)
	err = row.Scan(&user.ID, &user.Password)
	return
}
