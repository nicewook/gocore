package usecase_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nicewook/gocore/internal/config"
	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
	"github.com/nicewook/gocore/internal/usecase"
	"github.com/nicewook/gocore/pkg/security"
)

// 테스트용 RSA 키 생성
func generateTestRSAKeys() (string, string, error) {
	// 2048비트 RSA 키 생성
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// 개인키를 PEM 형식으로 인코딩
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 공개키를 PEM 형식으로 인코딩
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

func TestSignUpUser(t *testing.T) {
	// 테스트용 RSA 키 생성
	privateKeyPEM, publicKeyPEM, err := generateTestRSAKeys()
	if err != nil {
		t.Fatalf("Failed to generate RSA keys: %v", err)
	}

	// Create a test config
	cfg := &config.Config{
		Secure: config.SecureConfig{
			JWT: config.JWTConfig{
				PrivateKey:           privateKeyPEM,
				PublicKey:            publicKeyPEM,
				AccessExpirationMin:  60, // 60분
				RefreshExpirationDay: 30, // 30일
				Cookie: config.CookieConfig{
					Secure:   false,
					HTTPOnly: true,
					SameSite: "Lax",
					Domain:   "localhost",
				},
			},
		},
	}

	tests := []struct {
		name       string
		mockInput  *domain.User
		mockReturn *domain.User
		mockError  error
		expected   *domain.User
		expectErr  error
	}{
		{
			name: "Success",
			mockInput: &domain.User{
				Email:    "test@example.com",
				Password: "password123",
				Roles:    []string{domain.RoleUser},
			},
			mockReturn: &domain.User{
				ID:       1,
				Email:    "test@example.com",
				Password: "hashedpassword", // will be different in actual test
				Roles:    []string{domain.RoleUser},
			},
			mockError: nil,
			expected: &domain.User{
				ID:       1,
				Email:    "test@example.com",
				Password: "hashedpassword", // will be different in actual test
				Roles:    []string{domain.RoleUser},
			},
			expectErr: nil,
		},
		{
			name: "AlreadyExists",
			mockInput: &domain.User{
				Email:    "test@example.com",
				Password: "password123",
				Roles:    []string{domain.RoleUser},
			},
			mockReturn: nil,
			mockError:  domain.ErrAlreadyExists,
			expected:   nil,
			expectErr:  domain.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(mocks.AuthRepository)
			mockUserRepo := new(mocks.UserRepository)

			// We need to use a matcher for password since it will be hashed
			mockAuthRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(tt.mockReturn, tt.mockError)

			uc, err := usecase.NewAuthUseCase(mockAuthRepo, mockUserRepo, cfg)
			assert.NoError(t, err)

			ctx := context.Background()
			result, err := uc.SignUpUser(ctx, tt.mockInput)

			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Email, result.Email)
				assert.Equal(t, tt.expected.Roles, result.Roles)
			}

			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	privateKeyPEM, publicKeyPEM, err := generateTestRSAKeys()
	assert.NoError(t, err)

	cfg := &config.Config{
		Secure: config.SecureConfig{
			JWT: config.JWTConfig{
				PrivateKey:           privateKeyPEM,
				PublicKey:            publicKeyPEM,
				AccessExpirationMin:  15,
				RefreshExpirationDay: 7,
				Cookie: config.CookieConfig{
					Secure:   false,
					HTTPOnly: true,
					SameSite: "Lax",
					Domain:   "localhost",
				},
			},
		},
	}

	// 테스트 사용자 생성
	hashedPassword, err := security.GeneratePasswordHash("password", nil)
	assert.NoError(t, err)

	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
		Roles:    []string{domain.RoleUser},
	}

	tests := []struct {
		name      string
		email     string
		password  string
		mockUser  *domain.User
		mockError error
		expected  *domain.LoginResponse
		expectErr error
	}{
		{
			name:      "Success",
			email:     "test@example.com",
			password:  "password",
			mockUser:  user,
			mockError: nil,
			expected: &domain.LoginResponse{
				ID:    1,
				Email: "test@example.com",
			},
			expectErr: nil,
		},
		{
			name:      "User Not Found",
			email:     "nonexistent@example.com",
			password:  "password",
			mockUser:  nil,
			mockError: domain.ErrNotFound,
			expected:  nil,
			expectErr: domain.ErrNotFound,
		},
		{
			name:      "Invalid Password",
			email:     "test@example.com",
			password:  "wrongpassword",
			mockUser:  user,
			mockError: nil,
			expected:  nil,
			expectErr: errors.New("invalid credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(mocks.AuthRepository)
			mockUserRepo := new(mocks.UserRepository)

			if tt.email != "" {
				mockUserRepo.On("GetUserByEmail", mock.Anything, tt.email).Return(tt.mockUser, tt.mockError)
			}

			uc, err := usecase.NewAuthUseCase(mockAuthRepo, mockUserRepo, cfg)
			assert.NoError(t, err)

			ctx := context.Background()
			result, err := uc.Login(ctx, tt.email, tt.password)

			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	// 테스트용 RSA 키 생성
	privateKeyPEM, publicKeyPEM, err := generateTestRSAKeys()
	assert.NoError(t, err)

	// Create a test config
	cfg := &config.Config{
		Secure: config.SecureConfig{
			JWT: config.JWTConfig{
				PrivateKey:           privateKeyPEM,
				PublicKey:            publicKeyPEM,
				AccessExpirationMin:  60,
				RefreshExpirationDay: 30,
				Cookie: config.CookieConfig{
					Secure:   false,
					HTTPOnly: true,
					SameSite: "Lax",
					Domain:   "localhost",
				},
			},
		},
	}

	tests := []struct {
		name      string
		userID    int64
		expectErr error
	}{
		{
			name:      "Success",
			userID:    1,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(mocks.AuthRepository)
			mockUserRepo := new(mocks.UserRepository)

			uc, err := usecase.NewAuthUseCase(mockAuthRepo, mockUserRepo, cfg)
			assert.NoError(t, err)

			ctx := context.Background()
			err = uc.Logout(ctx, tt.userID)

			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	// 테스트용 RSA 키 생성
	privateKeyPEM, publicKeyPEM, err := generateTestRSAKeys()
	assert.NoError(t, err)

	// Create a test config
	cfg := &config.Config{
		Secure: config.SecureConfig{
			JWT: config.JWTConfig{
				PrivateKey:           privateKeyPEM,
				PublicKey:            publicKeyPEM,
				AccessExpirationMin:  15,
				RefreshExpirationDay: 7,
				Cookie: config.CookieConfig{
					Secure:   false,
					HTTPOnly: true,
					SameSite: "Lax",
					Domain:   "localhost",
				},
			},
		},
	}

	// 테스트 사용자 생성
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashedpassword",
		Roles:    []string{domain.RoleUser},
	}

	// 실제 리프레시 토큰 생성
	privateKey, _ := security.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	refreshTokenExpiration := time.Duration(cfg.Secure.JWT.RefreshExpirationDay) * 24 * time.Hour
	validRefreshToken, _ := security.GenerateRefreshToken(user.ID, user.Email, user.Roles, privateKey, refreshTokenExpiration)

	// 유효하지 않은 토큰
	invalidRefreshToken := "invalid.refresh.token"

	tests := []struct {
		name           string
		refreshToken   string
		mockUser       *domain.User
		mockError      error
		expectedResult *domain.LoginResponse
		expectErr      error
	}{
		{
			name:         "Success",
			refreshToken: validRefreshToken,
			mockUser:     user,
			mockError:    nil,
			expectedResult: &domain.LoginResponse{
				ID:    user.ID,
				Email: user.Email,
			},
			expectErr: nil,
		},
		{
			name:           "Invalid Token",
			refreshToken:   invalidRefreshToken,
			mockUser:       nil,
			mockError:      nil,
			expectedResult: nil,
			expectErr:      domain.ErrUnauthorized,
		},
		{
			name:           "User Not Found",
			refreshToken:   validRefreshToken,
			mockUser:       nil,
			mockError:      domain.ErrNotFound,
			expectedResult: nil,
			expectErr:      domain.ErrUnauthorized,
		},
		{
			name:           "Database Error",
			refreshToken:   validRefreshToken,
			mockUser:       nil,
			mockError:      domain.ErrInternal,
			expectedResult: nil,
			expectErr:      domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(mocks.AuthRepository)
			mockUserRepo := new(mocks.UserRepository)

			// 유효한 토큰이고 사용자 조회 모킹이 필요한 경우
			if tt.refreshToken == validRefreshToken && tt.name != "Invalid Token" {
				mockUserRepo.On("GetByID", mock.Anything, user.ID).Return(tt.mockUser, tt.mockError)
			}

			uc, err := usecase.NewAuthUseCase(mockAuthRepo, mockUserRepo, cfg)
			assert.NoError(t, err)

			ctx := context.Background()
			result, err := uc.RefreshToken(ctx, tt.refreshToken)

			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.Email, result.Email)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
