package service

import "mainstreet/internal/repository"

type Service struct{}

func New(repo *repository.Repository) *Service {
	return &Service{}
}
