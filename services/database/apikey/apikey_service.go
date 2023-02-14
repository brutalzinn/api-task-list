package apikey_service

import (
	"api-auto-assistant/db"
	entities "api-auto-assistant/models"
	"time"
)

func Delete(userId int64) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("DELETE FROM api_keys WHERE user_id=$1", userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Get(userId int64) (apikey entities.ApiKey, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT * FROM api_keys WHERE user_id=$1", userId)
	err = row.Scan(&apikey.ID, &apikey.ApiKey, &apikey.Scopes, &apikey.UserId, &apikey.Create_at, &apikey.Update_at)
	return
}
func GetAll(userId int64) (apikeys []entities.ApiKey, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT * FROM api_keys WHERE user_id=$1", userId)
	if err != nil {
		return
	}
	for rows.Next() {
		var apikey entities.ApiKey
		err = rows.Scan(&apikey.ID, &apikey.ApiKey, &apikey.Scopes, &apikey.UserId, &apikey.Create_at, &apikey.Update_at)
		if err != nil {
			continue
		}
		apikeys = append(apikeys, apikey)
	}
	return
}

func Count(userId int64) (count int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	row := conn.QueryRow("SELECT COUNT(*) FROM api_keys WHERE user_id=$1", userId)
	err = row.Scan(&count)
	return
}

func Insert(apiKey entities.ApiKey) (id int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()
	sql := "INSERT INTO api_keys (apikey, scopes, user_id, create_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err = conn.QueryRow(sql, &apiKey.ApiKey, &apiKey.Scopes, apiKey.UserId, time.Now()).Scan(&id)
	return
}
func Update(id int64, apikey entities.ApiKey) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	res, err := conn.Exec("UPDATE api_keys SET apikey=$1,scopes=$2,user_id=$3,update_at WHERE id=$4", &apikey.ApiKey, &apikey.Scopes, &apikey.UserId, time.Now(), id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
