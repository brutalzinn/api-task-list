package user_service

import (
	"time"

	"github.com/brutalzinn/api-task-list/db"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

func Delete(id string) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
func GetAll() (users []database_entities.User, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		var user database_entities.User
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
func Get(id int64) (user database_entities.User, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM users WHERE id=$1", id)
	err = row.Scan(&user.ID, &user.Email, &user.Password,
		&user.Username, &user.FirebaseToken, &user.CreateAt,
		&user.UpdateAt)
	return
}

func Insert(user database_entities.User) (id string, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	sql := "INSERT INTO users (email, password, username, firebaseToken, create_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = conn.QueryRow(ctx, sql, user.Email, user.Password, user.Username, user.FirebaseToken, time.Now()).Scan(&id)
	return
}
func Update(id int64, user database_entities.User) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	res, err := conn.Exec(ctx, "UPDATE users SET email=$1, password=$2, username=$3, firebaseToken=$4, update_at=$5 WHERE id=$6",
		user.Email, user.Password, user.Username, user.FirebaseToken, user.UpdateAt, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
func FindByEmail(email string) (user database_entities.User, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT id, password FROM users WHERE email=$1", email)
	err = row.Scan(&user.ID, &user.Password)
	return
}
