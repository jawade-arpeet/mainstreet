package middleware

import "mainstreet/internal/repository"

type Middleware struct {
	Request *RequestMiddleware
}

func New(repo *repository.Repository) *Middleware {
	return &Middleware{
		Request: newRequestMiddleware(),
	}
}
