package oauth_service

import (
	"log"

	"github.com/brutalzinn/api-task-list/db"
	database_entities "github.com/brutalzinn/api-task-list/models/database"
)

func Delete(oauthApp database_entities.OAuthApp) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "DELETE FROM users_oauth_client WHERE oauth_client_id=$1", oauthApp.OAuthClientId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func UpdateSecret(oauthApp database_entities.OAuthApp) (int64, error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)
	res, err := conn.Exec(ctx, "UPDATE FROM users_oauth_client WHERE oauth_client_id=$1", oauthApp.OAuthClientId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func List(userId string) (oauthApps []database_entities.OAuthApp, err error) {
	conn, err, ctx := db.OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "SELECT client_app.id, appname, mode, client_app.oauth_client_id, user_id, create_at, update_at FROM oauth_client_application as client_app INNER JOIN users_oauth_client as uc ON uc.oauth_client_id = client_app.oauth_client_id where uc.user_id=$1", userId)
	if err != nil {
		log.Fatal(err)
		return
	}
	for rows.Next() {
		var oauthApp database_entities.OAuthApp
		err = rows.Scan(&oauthApp.ID, &oauthApp.AppName, &oauthApp.Mode, &oauthApp.OAuthClientId, &oauthApp.UserId, &oauthApp.CreateAt, &oauthApp.UpdateAt)
		if err != nil {
			continue
		}
		oauthApps = append(oauthApps, oauthApp)
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
	_, err = tx.Exec(ctx, "INSERT INTO users_oauth_client (user_id, oauth_client_id) VALUES ($1, $2)", oAuthApp.UserId, oAuthApp.OAuthClientId)
	if err != nil {
		log.Fatal(err)
		return
	}
	nestedTx, err := tx.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, err = nestedTx.Exec(ctx, "INSERT INTO oauth_client_application (appname, mode, oauth_client_id) VALUES ($1, $2, $3)",
		oAuthApp.AppName, oAuthApp.Mode, oAuthApp.OAuthClientId)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}
	return
}
