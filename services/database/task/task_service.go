package task_service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brutalzinn/api-task-list/db"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

func Delete(id int64) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func DeleteTasksByRepo(repo_id int64) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM tasks WHERE repo_id=$1", repo_id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
func ReplaceTasksByRepo(tasks []database_entities.Task, repo_id int64) (err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec(ctx, "DELETE FROM tasks WHERE repo_id=$1", repo_id)
	if err != nil {
		log.Fatal(err)
	}
	nestedTx, err := tx.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, task := range tasks {
		_, err = nestedTx.Exec(ctx, "INSERT INTO tasks (title, repo_id, description, text, create_at) VALUES ($1, $2, $3, $4, $5)", task.Title, repo_id, task.Description, task.Text, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return
}
func UpdateTasks(tasks []database_entities.Task) (err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(ctx)

	if err != nil {
		log.Fatal(err)
	}
	for _, task := range tasks {
		_, err = tx.Exec(ctx, "UPDATE tasks SET title=$1,repo_id=$2,description=$3,text=$4, update_at=$5 WHERE id=$6", task.Title, task.RepoId, task.Description, task.Text, time.Now(), task.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return
}
func GetAll() (tasks []database_entities.Task, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := conn.Query(ctx, "SELECT * FROM tasks")
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
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM tasks WHERE id=$1", id)
	err = row.Scan(&task.ID, &task.Title, &task.RepoId, &task.Description, &task.Text, &task.CreateAt, &task.UpdateAt)
	return
}

func Insert(task database_entities.Task) (id int64, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	sql := "INSERT INTO tasks (title, repo_id, description, text, create_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = conn.QueryRow(ctx, sql, &task.Title, &task.RepoId, &task.Description, &task.Text, time.Now()).Scan(&id)
	return
}
func Update(id int64, task database_entities.Task) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "UPDATE tasks SET title=$1,description=$2,text=$3, update_at=$4 WHERE id=$5", &task.Title, &task.Description, &task.Text, time.Now(), id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func Paginate(repo_id int64, limit int64, offset int64, order string) (tasks []database_entities.Task, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	query := fmt.Sprintf("SELECT * FROM tasks WHERE repo_id=$1 ORDER BY create_at %s LIMIT $2 OFFSET $3", order)
	rows, err := conn.Query(ctx, query, repo_id, limit, offset)
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
func Count(repo_id int64) (count int64, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT COUNT(*) FROM tasks where repo_id=$1", repo_id)
	err = row.Scan(&count)
	return
}
