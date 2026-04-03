package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustConnect(dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("db connect error: ", err)
	}
	return pool
}
