package oauth_service

import (
	"github.com/brutalzinn/api-task-list/db"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

func List(userId string) (apiKeys []database_entities.ApiKey, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "SELECT * FROM oauth2_clients WHERE user_id=$1", userId)
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
