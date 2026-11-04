package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
	"github.com/Project-TETI/AgriSense/backend/internal/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(context.Context, domain.RegisterRequest) (*domain.AuthSession, error)
	Login(context.Context, domain.LoginRequest) (*domain.AuthSession, error)
	RefreshToken(context.Context, string) (*domain.AuthSession, error)
	Logout(context.Context, string) error
	GetProfile(context.Context, uuid.UUID) (*domain.UserResponse, error)
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	transaction      repository.AuthTransaction
	jwtSecret        string
	accessDuration   time.Duration
	refreshDuration  time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	transaction repository.AuthTransaction,
	jwtSecret string,
	accessDuration time.Duration,
	refreshDuration time.Duration,
) AuthService {
	return &authService{
		userRepo: userRepo, refreshTokenRepo: refreshTokenRepo, transaction: transaction,
		jwtSecret: jwtSecret, accessDuration: accessDuration, refreshDuration: refreshDuration,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthSession, error) {
	if len(req.Password) > 72 {
		return nil, domain.ErrInvalidInput
	}
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find registration email: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyRegistered
	}
	// Hashing is deliberately outside the transaction to avoid holding locks during bcrypt.
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash registration password: %w", err)
	}
	user := domain.User{Name: req.Name, Email: req.Email, PasswordHash: passwordHash}
	var session *domain.AuthSession
	err = s.transaction.Run(ctx, func(users repository.UserRepository, tokens repository.RefreshTokenRepository) error {
		if err := users.Create(ctx, &user); err != nil {
			return err
		}
		var err error
		session, err = s.generateAuthSession(ctx, tokens, &user, time.Now().UTC().Add(s.refreshDuration))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("register session: %w", err)
	}
	return session, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthSession, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find login user: %w", err)
	}
	if user == nil || !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}
	return s.generateAuthSession(ctx, s.refreshTokenRepo, user, time.Now().UTC().Add(s.refreshDuration))
}

func validRefreshToken(raw string) bool {
	if len(raw) != 64 {
		return false
	}
	_, err := hex.DecodeString(raw)
	return err == nil
}

func (s *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (*domain.AuthSession, error) {
	if !validRefreshToken(rawRefreshToken) {
		return nil, domain.ErrInvalidRefreshToken
	}
	var session *domain.AuthSession
	err := s.transaction.Run(ctx, func(users repository.UserRepository, tokens repository.RefreshTokenRepository) error {
		storedToken, err := tokens.FindByTokenHashForUpdate(ctx, utils.HashToken(rawRefreshToken))
		if err != nil {
			return err
		}
		// Check time after acquiring the row lock, not before waiting for it.
		if storedToken == nil || storedToken.RevokedAt != nil || !storedToken.ExpiresAt.After(time.Now().UTC()) {
			return domain.ErrInvalidRefreshToken
		}
		user, err := users.FindByID(ctx, storedToken.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return domain.ErrInvalidRefreshToken
		}
		revoked, err := tokens.Revoke(ctx, storedToken.ID)
		if err != nil {
			return err
		}
		if !revoked {
			return domain.ErrInvalidRefreshToken
		}
		session, err = s.generateAuthSession(ctx, tokens, user, storedToken.ExpiresAt)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("rotate refresh token: %w", err)
	}
	return session, nil
}

func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	if !validRefreshToken(rawRefreshToken) {
		return nil
	}
	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, utils.HashToken(rawRefreshToken))
	if err != nil {
		return fmt.Errorf("find logout token: %w", err)
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return nil
	}
	_, err = s.refreshTokenRepo.Revoke(ctx, storedToken.ID)
	if err != nil {
		return fmt.Errorf("revoke logout token: %w", err)
	}
	return nil
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find profile: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return &domain.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, CreatedAt: user.CreatedAt}, nil
}

func (s *authService) generateAuthSession(ctx context.Context, tokens repository.RefreshTokenRepository, user *domain.User, refreshExpiresAt time.Time) (*domain.AuthSession, error) {
	now := time.Now().UTC().Truncate(time.Second)
	accessExpiresAt := now.Add(s.accessDuration).Truncate(time.Second)
	if refreshExpiresAt.Before(accessExpiresAt) {
		accessExpiresAt = refreshExpiresAt.Truncate(time.Second)
	}
	if !accessExpiresAt.After(now) {
		return nil, domain.ErrInvalidRefreshToken
	}
	accessToken, err := utils.GenerateAccessTokenWithExpiry(user.ID, user.Email, user.Name, s.jwtSecret, now, accessExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}
	rawRefreshToken, tokenHash, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshToken := domain.RefreshToken{UserID: user.ID, TokenHash: tokenHash, ExpiresAt: refreshExpiresAt}
	if err := tokens.Create(ctx, &refreshToken); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}
	return &domain.AuthSession{
		Response: domain.AuthResponse{
			AccessToken: accessToken, TokenType: "Bearer", ExpiresIn: int64(accessExpiresAt.Sub(now).Seconds()),
			User: domain.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, CreatedAt: user.CreatedAt},
		},
		RefreshToken: rawRefreshToken, RefreshExpiresAt: refreshExpiresAt,
	}, nil
}
