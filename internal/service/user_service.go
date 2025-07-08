package service

import (
	"context"
	"demo/internal/repository"
	"demo/pkg/logger"

	"go.uber.org/zap"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) (map[string]string, error) {
	logger := logger.FromContext(ctx)
	logger.Info("Service: GetUser called", zap.String("id", id))

	return s.repo.FindUser(ctx, id)
}
