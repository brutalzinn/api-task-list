package db

import (
	"context"
	"fmt"

	"github.com/brutalzinn/api-task-list/configs"
	"github.com/jackc/pgx/v4"
)

func OpenConnection() (*pgx.Conn, error, context.Context) {
	ctx := context.Background()
	conn, err := pgx.ConnectConfig(ctx, GetConnectionAdapter())
	if err != nil {
		panic(err)
	}
	err = conn.Ping(ctx)
	return conn, err, ctx
}

func GetConnectionAdapter() *pgx.ConnConfig {
	conf := configs.GetConfig().DB
	sc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Pass, conf.Database)
	pgConfig, _ := pgx.ParseConfig(sc)
	return pgConfig
}
