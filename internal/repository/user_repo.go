package repository

import (
	"context"
	"demo/pkg/logger"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindUser(ctx context.Context, id string) (map[string]string, error) {
	logger := logger.FromContext(ctx)
	tracer := trace.SpanFromContext(ctx).TracerProvider().Tracer("repository")
	ctx, span := tracer.Start(ctx, "Repository:FindUser")
	defer span.End()

	logger.Info("Repository: FindUser called", zap.String("id", id))

	// Mock DB
	user := map[string]string{
		"id":   id,
		"name": "John Doe",
	}
	return user, nil
}
