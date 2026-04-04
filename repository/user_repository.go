package repository

import (
	"blog/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUserName(ctx context.Context, username string) (models.Users, error) {
	var user models.Users
	err := r.db.QueryRow(ctx, `SELECT * FROM users WHERE username=$1`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		err = errors.New("user not found")
		return user, err
	}
	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (models.MeResponse, error) {
	var req models.MeResponse
	err := r.db.QueryRow(ctx, `SELECT id, username, email FROM users WHERE id=$1`, id).Scan(&req.ID, &req.Username, &req.Email)
	return req, err
}

func (r *UserRepository) Add(ctx context.Context, username string, email string, hPassword string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
	`, username, email, hPassword)
	return err
}
