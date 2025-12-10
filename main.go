package main

import (
	"blog/router"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
		return
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is empty")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connected DB: ", err)
	}
	defer pool.Close()

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS posts (
					id         SERIAL PRIMARY KEY,
					title      TEXT        NOT NULL,
					content    TEXT        NOT NULL,
					category   TEXT,
					tags TEXT[],
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
				);
			`)
	if err != nil {
		log.Fatal("Create table error:", err)
	}
	log.Println("Table posts is ready")

	r := router.SetupRouter()
	r.Run(":8080")
}
