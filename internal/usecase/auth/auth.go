package auth

import (
	"bukus/internal/entity"
	"bukus/internal/usecase"
	"bukus/pkg/jwt"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

// UseCase -.
type UseCase struct {
	repo       usecase.UserRepo
	jwtAccess  *jwt.Manager
	jwtRefresh *jwt.Manager
}

// New -.
func New(repo usecase.UserRepo, jwtAccess, jwtRefresh *jwt.Manager) *UseCase {
	return &UseCase{
		repo:       repo,
		jwtAccess:  jwtAccess,
		jwtRefresh: jwtRefresh,
	}
}

// Register -.
func (obj *UseCase) Register(ctx context.Context, email, password string) error {
	// Check if user exists
	_, err := obj.repo.FindByEmail(ctx, email)
	if err == nil {
		return ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("auth usecase - Register - bcrypt: %w", err)
	}

	user := entity.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := obj.repo.Store(ctx, user); err != nil {
		return fmt.Errorf("auth usecase - Register - Store: %w", err)
	}

	return nil
}

// Login -.
func (obj *UseCase) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := obj.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	return obj.generateTokenPair(user.ID)
}

// Refresh -.
func (obj *UseCase) Refresh(_ context.Context, refreshToken string) (string, string, error) {
	userID, err := obj.jwtRefresh.ParseToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	return obj.generateTokenPair(userID)
}

func (obj *UseCase) generateTokenPair(userID string) (string, string, error) {
	accessToken, err := obj.jwtAccess.GenerateToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("auth usecase - generateTokenPair - access: %w", err)
	}

	refreshToken, err := obj.jwtRefresh.GenerateToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("auth usecase - generateTokenPair - refresh: %w", err)
	}

	return accessToken, refreshToken, nil
}
