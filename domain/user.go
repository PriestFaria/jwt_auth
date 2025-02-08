package domain

import "errors"

var (
	UserNotFoundError       = errors.New("user not found")
	InvalidCredentialsError = errors.New("invalid credentials")
)

type User struct {
	ID       int64
	Email    string
	Password string
	Role     Role
}

func (u User) isAdmin() bool {
	return u.Role == AdminRole
}
