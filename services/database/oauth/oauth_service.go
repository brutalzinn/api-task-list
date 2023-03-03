package oauth_service

import (
	"log"

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
func CreateOauthForUser(oAuthApp database_entities.OAuthApp) (err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec(ctx, "INSERT INTO users_oauth_client (user_id, oauth_client_id) VALUES($1, $2)", oAuthApp.UserId, oAuthApp.OAuthClientId)
	if err != nil {
		log.Fatal(err)
	}
	nestedTx, err := tx.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, err = nestedTx.Exec(ctx, "INSERT INTO oauth_client_application(appname, mode, oauth_client_id) VALUES($1, $2, $3)",
		oAuthApp.AppName, oAuthApp.Mode, oAuthApp.OAuthClientId, oAuthApp.CreateAt, oAuthApp.UpdateAt)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return
}
