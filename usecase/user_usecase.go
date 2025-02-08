package usecase

//TODO разделить реаолизации и интерфейсы — щас в коде блуждать можно просто
import (
	"auth_backend/domain"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}
type TokenGenerator interface {
	GenerateToken(user *domain.User) (string, error)
}

type userUseCase struct {
	userRepository UserRepository
	tokenGenerator TokenGenerator
}

type UserUseCase interface {
	Authenticate(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password string) error
}

func NewUserUseCase(repository UserRepository, tokenGenerator TokenGenerator) *userUseCase {
	return &userUseCase{
		userRepository: repository,
		tokenGenerator: tokenGenerator,
	}
}

func (u *userUseCase) Authenticate(ctx context.Context, email string, password string) (string, error) {
	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", domain.UserNotFoundError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.InvalidCredentialsError
	}

	token, err := u.tokenGenerator.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUseCase) Register(ctx context.Context, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &domain.User{
		Email:    email,
		Password: string(hashed),
	}

	err = u.userRepository.CreateUser(ctx, user)

	if err != nil {
		return err
	}

	return nil
}
