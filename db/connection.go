package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func GetConnectionUri() string {
	uri := os.Getenv("PG_URI")
	fmt.Println(uri)
	if uri == "" {
		fmt.Println("Env variable PG_URI is required to run")
		os.Exit(1)
	}
	return uri
}
func OpenConnection() (*pgx.Conn, error, context.Context) {
	ctx := context.TODO()
	conn, err := pgx.Connect(ctx, os.Getenv("PG_URI"))
	if err != nil {
		panic(err)
	}
	err = conn.Ping(ctx)
	return conn, err, ctx
}
