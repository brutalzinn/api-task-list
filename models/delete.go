package models

import "api-auto-assistant/db"

func Delete(id int64, task Task) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("DELETE tasks WHERE id=$1", task.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
