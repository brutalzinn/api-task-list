package models

import "api-auto-assistant/db"

func Update(id int64, task Task) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("UPDATE tasks SET title=$1, description=$2 WHERE id=$3", task.Title, task.Description, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
