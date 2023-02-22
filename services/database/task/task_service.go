package task_service

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
	res, err := conn.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func GetAll() (tasks []database_entities.Task, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT * FROM tasks")
	if err != nil {
		return
	}
	for rows.Next() {
		var task database_entities.Task
		err = rows.Scan(&task.ID, &task.Title, &task.RepoId, &task.Description, &task.CreateAt, &task.UpdateAt)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return
}
func Get(id int64) (task database_entities.Task, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT * FROM tasks WHERE id=$1", id)
	err = row.Scan(&task.ID, &task.Title, &task.RepoId, &task.Description, &task.CreateAt, &task.UpdateAt)
	return
}

func Insert(task database_entities.Task) (id int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	sql := "INSERT INTO tasks (title, repo_id, description, create_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err = conn.QueryRow(sql, &task.Title, &task.RepoId, &task.Description, time.Now()).Scan(&id)
	return
}
func Update(id int64, task database_entities.Task) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("UPDATE tasks SET title=$1,description=$2,update_at=$3 WHERE id=$4", &task.Title, &task.Description, time.Now(), id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func Paginate(repo_id int, limit int, offset int, order string) (tasks []database_entities.Task, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	query := fmt.Sprintf("SELECT * FROM tasks WHERE repo_id=$1 ORDER BY create_at %s LIMIT $2 OFFSET $3", order)
	rows, err := conn.Query(query, repo_id, limit, offset)
	if err != nil {
		return
	}
	for rows.Next() {
		var task database_entities.Task
		err = rows.Scan(&task.ID, &task.Title, &task.RepoId, &task.Description, &task.CreateAt, &task.UpdateAt)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return
}
func Count() (count int, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT COUNT(*) FROM tasks")
	err = row.Scan(&count)
	return
}
