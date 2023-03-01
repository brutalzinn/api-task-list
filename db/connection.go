package db

import (
	"database/sql"
	"fmt"

	"github.com/brutalzinn/api-task-list/configs"
	"github.com/jackc/pgx/v4"
)

func OpenConnection() (*sql.DB, error) {
	conf := configs.GetConfig().DB
	sc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Pass, conf.Database)
	conn, err := sql.Open("postgres", sc)
	if err != nil {
		panic(err)
	}
	err = conn.Ping()
	return conn, err
}

func GetConnectionAdapter() *pgx.ConnConfig {
	conf := configs.GetConfig().DB
	sc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Pass, conf.Database)
	pgConfig, _ := pgx.ParseConfig(sc)
	return pgConfig
}
