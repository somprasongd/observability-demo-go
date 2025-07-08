package repository

import (
	"context"
	"log"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindUser(ctx context.Context, id string) (map[string]string, error) {
	log.Printf("Repository: FindUser called with ID %s\n", id)

	// Mock DB
	user := map[string]string{
		"id":   id,
		"name": "John Doe",
	}
	return user, nil
}
