package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var uri string

func CreateConnection() {
	uri := os.Getenv("PG_URI")
	if uri == "" {
		fmt.Println("Env variable PG_URI is required to run")
		os.Exit(1)
	}
}
func GetConnectionUri() string {
	return uri
}
func OpenConnection() (*pgx.Conn, error, context.Context) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, uri)
	if err != nil {
		panic(err)
	}
	err = conn.Ping(ctx)
	return conn, err, ctx
}
