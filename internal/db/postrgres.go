package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres() (*pgxpool.Pool, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}