package service

import (
	"context"
	"demo/internal/repository"
	"log"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) (map[string]string, error) {
	log.Printf("Service: GetUser called with ID %s\n", id)
	return s.repo.FindUser(ctx, id)
}
