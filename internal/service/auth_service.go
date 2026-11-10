package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/model"
	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/hash"
	"github.com/amarnathdbg101-coder/GoProjectStructure/pkg/jwt"
	"github.com/amarnathdbg101-coder/GoProjectStructure/internal/repository"
)

// AuthService coordinates authentication business logic.
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// NewAuthService constructs an AuthService with required dependencies.
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account with a hashed password and returns an auth token.
func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Role:         "user",
		IsActive:     true,
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(created.ID, created.Email, created.Role, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User:  created,
	}, nil
}

// Login verifies user credentials and returns a JWT auth response.
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated; please contact support")
	}

	if !hash.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetProfile retrieves the profile for an authenticated user.
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
