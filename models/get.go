package models

import "api-auto-assistant/db"

func Get(id int64) (task Task, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT * FROM Tasks WHERE id=$1", id)
	err = row.Scan(&task.ID, &task.Title, &task.Description)
	return
}
