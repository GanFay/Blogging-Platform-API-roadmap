package handlers

import (
	"blog/repository"
)

type Handler struct {
	Posts *repository.PostRepository
	Users *repository.UserRepository
}

func NewHandler(p *repository.PostRepository, u *repository.UserRepository) *Handler {
	return &Handler{Posts: p, Users: u}
}
