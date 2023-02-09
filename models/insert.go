package models

import "api-auto-assistant/db"

func Insert(task Task) (id int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	sql := "INSERT INTO Tasks (title, description) VALUES ($1, $2) RETURNING id"
	err = conn.QueryRow(sql, task.Title, task.Description).Scan(&id)
	return
}
