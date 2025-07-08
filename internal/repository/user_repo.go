package repository

import (
	"context"
	"demo/pkg/logger"

	"go.uber.org/zap"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindUser(ctx context.Context, id string) (map[string]string, error) {
	logger := logger.FromContext(ctx)
	logger.Info("Repository: FindUser called", zap.String("id", id))

	// Mock DB
	user := map[string]string{
		"id":   id,
		"name": "John Doe",
	}
	return user, nil
}
