package apikey_service

import (
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
	res, err := conn.Exec(ctx, "DELETE FROM api_keys WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func Get(keyId string) (apiKey database_entities.ApiKey, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM api_keys WHERE id=$1", keyId)
	err = row.Scan(&apiKey.ID, &apiKey.ApiKey, &apiKey.Scopes, &apiKey.UserId, &apiKey.Name, &apiKey.NameNormalized, &apiKey.ExpireAt, &apiKey.CreateAt, &apiKey.UpdateAt)
	return
}
func GetAll(userId string) (apiKeys []database_entities.ApiKey, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "SELECT * FROM api_keys WHERE user_id=$1", userId)
	if err != nil {
		return
	}
	for rows.Next() {
		var apiKey database_entities.ApiKey
		err = rows.Scan(&apiKey.ID, &apiKey.ApiKey, &apiKey.Scopes, &apiKey.UserId, &apiKey.Name, &apiKey.NameNormalized, &apiKey.ExpireAt, &apiKey.CreateAt, &apiKey.UpdateAt)
		if err != nil {
			continue
		}
		apiKeys = append(apiKeys, apiKey)
	}
	return
}
func GetByUserAndName(userId int64, appName string) (apiKey database_entities.ApiKey, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM api_keys WHERE user_id=$1 and name_normalized=$2", userId, appName)
	err = row.Scan(&apiKey.ID, &apiKey.ApiKey, &apiKey.Scopes, &apiKey.UserId, &apiKey.Name, &apiKey.NameNormalized, &apiKey.ExpireAt, &apiKey.CreateAt, &apiKey.UpdateAt)
	return
}
func GetByIdAndUser(id string, userId string) (apiKey database_entities.ApiKey, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT * FROM api_keys WHERE id=$1 and user_id=$2", id, userId)
	err = row.Scan(&apiKey.ID, &apiKey.ApiKey, &apiKey.Scopes, &apiKey.UserId, &apiKey.Name, &apiKey.NameNormalized, &apiKey.ExpireAt, &apiKey.CreateAt, &apiKey.UpdateAt)
	return
}
func CountByUserAndName(userId string, appName string) (count int64, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT COUNT(*) FROM api_keys WHERE user_id=$1 and name_normalized=$2", userId, appName)
	err = row.Scan(&count)
	return
}

func DeleteByIdAndUser(id string, userId string) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM api_keys WHERE id=$1 and user_id=$2", id, userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
func Count(userId string) (count int64, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	row := conn.QueryRow(ctx, "SELECT COUNT(*) FROM api_keys WHERE user_id=$1", userId)
	err = row.Scan(&count)
	return
}
func Insert(apiKey database_entities.ApiKey) (id string, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	sql := "INSERT INTO api_keys (id, apiKey, scopes, name, name_normalized, user_id, create_at, expire_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err = conn.QueryRow(ctx, sql, &apiKey.ID, &apiKey.ApiKey, &apiKey.Scopes, &apiKey.Name, &apiKey.NameNormalized, apiKey.UserId, time.Now(), &apiKey.ExpireAt).Scan(&id)
	return
}
func Update(id string, apiKey database_entities.ApiKey) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "UPDATE api_keys SET apiKey=$1,scopes=$2,name=$3,update_at=$4,expire_at=$5 WHERE id=$6", &apiKey.ApiKey, &apiKey.Scopes, &apiKey.Name, time.Now(), &apiKey.ExpireAt, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
