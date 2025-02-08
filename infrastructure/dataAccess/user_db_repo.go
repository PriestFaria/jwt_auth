package dataAccess

import (
	"auth_backend/domain"
	"context"
	"database/sql"
	"fmt"
)

type UserRepositoryDB struct {
	db *sql.DB
}

func NewUserRepositoryDB(db *sql.DB) *UserRepositoryDB {
	return &UserRepositoryDB{db: db}
}

func (u *UserRepositoryDB) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password, role FROM users WHERE email = $1 LIMIT 1`
	row := u.db.QueryRowContext(ctx, query, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.UserNotFoundError
		}
		return nil, fmt.Errorf("Get by email: %w", err)
	}

	return &user, nil
}

func (u *UserRepositoryDB) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id`
	err := u.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
