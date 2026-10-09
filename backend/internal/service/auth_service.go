package service

import (
	"context"
	"errors"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
	"github.com/Project-TETI/AgriSense/backend/internal/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, string, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, string, error)
	RefreshToken(ctx context.Context, rawRefreshToken string) (*domain.AuthResponse, string, error)
	Logout(ctx context.Context, rawRefreshToken string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error)
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtSecret        string
	accessDuration   time.Duration
	refreshDuration  time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtSecret string,
	accessDuration time.Duration,
	refreshDuration time.Duration,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
		accessDuration:   accessDuration,
		refreshDuration:  refreshDuration,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, string, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", errors.New("email sudah terdaftar")
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	user := domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, "", err
	}

	return s.generateAuthSession(ctx, &user)
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("email atau kata sandi tidak valid")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, "", errors.New("email atau kata sandi tidak valid")
	}

	return s.generateAuthSession(ctx, user)
}

func (s *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (*domain.AuthResponse, string, error) {
	if rawRefreshToken == "" {
		return nil, "", errors.New("refresh token tidak ditemukan")
	}

	tokenHash := utils.HashToken(rawRefreshToken)
	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, "", err
	}
	if storedToken == nil {
		return nil, "", errors.New("refresh token tidak valid")
	}

	if storedToken.RevokedAt != nil {
		return nil, "", errors.New("refresh token sudah dicabut")
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, "", errors.New("refresh token telah kedaluwarsa")
	}

	if err := s.refreshTokenRepo.Revoke(ctx, storedToken.ID); err != nil {
		return nil, "", err
	}

	return s.generateAuthSession(ctx, &storedToken.User)
}

func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}

	tokenHash := utils.HashToken(rawRefreshToken)
	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}
	if storedToken != nil && storedToken.RevokedAt == nil {
		return s.refreshTokenRepo.Revoke(ctx, storedToken.ID)
	}

	return nil
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	return &domain.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) generateAuthSession(ctx context.Context, user *domain.User) (*domain.AuthResponse, string, error) {
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Name, s.jwtSecret, s.accessDuration)
	if err != nil {
		return nil, "", err
	}

	rawRefreshToken, tokenHash, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	refreshTokenModel := domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.refreshDuration),
	}
	if err := s.refreshTokenRepo.Create(ctx, &refreshTokenModel); err != nil {
		return nil, "", err
	}

	response := &domain.AuthResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.accessDuration.Seconds()),
		User: domain.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	return response, rawRefreshToken, nil
}
