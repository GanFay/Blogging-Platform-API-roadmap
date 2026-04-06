package main

import (
	"blog/db"
	"blog/handlers"
	"blog/repository"
	"blog/router"
	"log"
	"os"

	_ "blog/docs"

	"github.com/joho/godotenv"
)

// @title Blogging Platform API
// @version 1.0
// @description REST API for a blogging platform built with Go, Gin, PostgreSQL, JWT auth, refresh tokens and ownership protection.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is empty")
	}

	pool := db.MustConnect(dbURL)
	defer pool.Close()

	postRep := repository.NewPostRepository(pool)
	userRep := repository.NewUserRepository(pool)

	h := handlers.NewHandler(postRep, userRep)
	r := router.SetupRouter(h)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
