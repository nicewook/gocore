package usecase

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/nicewook/gocore/internal/config"
	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/pkg/security"
)

type authUseCase struct {
	authRepo   domain.AuthRepository
	userRepo   domain.UserRepository
	config     *config.Config
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewAuthUseCase(authRepo domain.AuthRepository, userRepo domain.UserRepository, config *config.Config) (domain.AuthUseCase, error) {
	// RSA 키 파싱
	privateKey, err := security.ParseRSAPrivateKeyFromPEM(config.Secure.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := security.ParseRSAPublicKeyFromPEM(config.Secure.JWT.PublicKey)
	if err != nil {
		return nil, err
	}

	return &authUseCase{
		authRepo:   authRepo,
		userRepo:   userRepo,
		config:     config,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (uc *authUseCase) SignUpUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// 비밀번호 해싱
	hashedPassword, err := security.GeneratePasswordHash(user.Password, nil)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return uc.authRepo.CreateUser(ctx, user)
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	// 이메일로 사용자 조회
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// 비밀번호 검증
	match, err := security.ComparePasswordHash(password, user.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("invalid credentials")
	}

	// 토큰 생성
	return uc.generateTokens(user)
}

// Logout 사용자 로그아웃 처리
func (uc *authUseCase) Logout(ctx context.Context, userID int64) error {
	// 현재 구현에서는 클라이언트 측에서 토큰을 삭제하는 방식으로 처리
	// 서버 측에서는 특별한 작업이 필요 없음

	// TODO: 향후 토큰 블랙리스트 구현 시 여기에 추가
	// - 사용자의 현재 토큰을 블랙리스트에 추가
	// - Redis 또는 다른 인메모리 저장소 사용 권장
	// - 토큰의 만료 시간까지만 블랙리스트에 유지

	return nil
}

// RefreshToken validates a refresh token and issues new access and refresh tokens
func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
	// Validate the refresh token
	claims, err := security.ValidateRefreshToken(refreshToken, uc.publicKey)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	// Get user by ID
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, err
	}

	response, err := uc.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (uc *authUseCase) generateTokens(user *domain.User) (*domain.LoginResponse, error) {
	// Generate access token
	accessTokenExpiration := time.Duration(uc.config.Secure.JWT.AccessExpirationMin) * time.Minute
	accessToken, err := security.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Roles,
		uc.privateKey,
		accessTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenExpiration := time.Duration(uc.config.Secure.JWT.RefreshExpirationDay) * 24 * time.Hour
	refreshToken, err := security.GenerateRefreshToken(
		user.ID,
		user.Email,
		user.Roles,
		uc.privateKey,
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	// Create response
	response := &domain.LoginResponse{
		ID:                     user.ID,
		Email:                  user.Email,
		AccessToken:            accessToken,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: time.Now().Add(refreshTokenExpiration),
	}
	return response, nil
}
