package middleware

import "mainstreet/internal/repository"

type Middleware struct{}

func New(repo *repository.Repository) *Middleware {
	return &Middleware{}
}
