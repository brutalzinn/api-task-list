package repo_service

import (
	"fmt"
	"time"

	"github.com/brutalzinn/api-task-list/db"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

func Delete(id int64) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("DELETE FROM repos WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func GetAll() (repos []database_entities.Repo, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT * FROM repos")
	if err != nil {
		return
	}
	for rows.Next() {
		var repo database_entities.Repo
		err = rows.Scan(&repo.ID, &repo.Title, &repo.Description, &repo.UserId, &repo.CreateAt, &repo.UpdateAt)
		if err != nil {
			continue
		}
		repos = append(repos, repo)
	}
	return
}
func Paginate(limit int, offset int, order string) (repos []database_entities.Repo, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	query := fmt.Sprintf("SELECT * FROM repos ORDER BY create_at %s LIMIT $1 OFFSET $2", order)
	rows, err := conn.Query(query, limit, offset)
	if err != nil {
		return
	}
	for rows.Next() {
		var repo database_entities.Repo
		err = rows.Scan(&repo.ID, &repo.Title, &repo.Description, &repo.UserId, &repo.CreateAt, &repo.UpdateAt)
		if err != nil {
			continue
		}
		repos = append(repos, repo)
	}
	return
}
func Count() (count int, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT COUNT(*) FROM repos")
	err = row.Scan(&count)
	return
}
func Get(id int64) (repo database_entities.Repo, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT * FROM repos WHERE id=$1", id)
	err = row.Scan(&repo.ID, &repo.Title, &repo.Description, &repo.UserId, &repo.CreateAt, &repo.UpdateAt)
	return
}

func Insert(repo database_entities.Repo) (id int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	sql := "INSERT INTO repos (title, description, user_id, create_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err = conn.QueryRow(sql, &repo.Title, &repo.Description, &repo.UserId, time.Now()).Scan(&id)
	return
}
func Update(id int64, repo database_entities.Repo) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("UPDATE repos SET title=$1,description=$2,user_id=$3,update_at=$4 WHERE id=$5", &repo.Title, &repo.Description, &repo.UserId, time.Now(), id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
