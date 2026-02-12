package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ai-companion-be/internal/config"
	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

// AuthService handles user authentication.
type AuthService struct {
	users  repository.UserRepository
	jwtCfg config.JWTConfig
}

// NewAuthService creates a new AuthService.
func NewAuthService(users repository.UserRepository, jwtCfg config.JWTConfig) *AuthService {
	return &AuthService{users: users, jwtCfg: jwtCfg}
}

// Signup registers a new user and returns an auth token.
func (s *AuthService) Signup(ctx context.Context, req models.SignupRequest) (*models.AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, fmt.Errorf("email, password, and name are required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashed),
		Name:     req.Name,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

// Login authenticates a user and returns an auth token.
func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

// GetUser returns a user by ID.
func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(s.jwtCfg.Expiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}
