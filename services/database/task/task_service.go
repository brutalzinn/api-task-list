package task_service

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
	res, err := conn.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func GetAll() (tasks []entities.Task, err error) {
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
		var task entities.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Create_at, &task.Update_at)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return
}
func Get(id int64) (task entities.Task, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT * FROM tasks WHERE id=$1", id)
	err = row.Scan(&task.ID, &task.Title, &task.Description, &task.Create_at, &task.Update_at)
	return
}

func Insert(task entities.Task) (id int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	sql := "INSERT INTO tasks (title, description, create_at) VALUES ($1, $2, $3) RETURNING id"
	err = conn.QueryRow(sql, &task.Title, &task.Description, time.Now()).Scan(&id)
	return
}
func Update(id int64, task entities.Task) (int64, error) {
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
