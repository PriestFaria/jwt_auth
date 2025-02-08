package dataAccess

import (
	"auth_backend/domain"
	"context"
)

type UserRepositoryMemory struct {
	users map[string]*domain.User
}

func NewUserRepositoryMemory() *UserRepositoryMemory {
	return &UserRepositoryMemory{
		users: map[string]*domain.User{
			"test@example.com": &domain.User{
				ID:       1,
				Email:    "test@example.com",
				Password: "secret",
			},
		},
	}
}

func (urm *UserRepositoryMemory) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, ok := urm.users[email]
	if !ok {
		return nil, domain.UserNotFoundError
	}
	return user, nil
}
