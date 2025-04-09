package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Validate는 기본 validator 이후에 실행될 커스텀 유효성 검사를 수행합니다
func (r *SignUpRequest) Validate(c echo.Context) error {
	// 기본 validation 태그 검증 (required, email, min 등)
	if err := c.Validate(r); err != nil {
		return err
	}

	// 추가적인 email 도메인 검증. hotmail.com 도메인은 허용하지 않음
	if strings.HasSuffix(r.Email, "@hotmail.com") {
		return errors.New("restricted email domain")
	}

	return nil
}

type SignUpResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID                     int64     `json:"id"`
	Email                  string    `json:"email"`
	AccessToken            string    `json:"access_token"`
	RefreshToken           string    `json:"-"` // Not included in JSON response
	RefreshTokenExpiration time.Time `json:"-"` // Not included in JSON response
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
}

type AuthUseCase interface {
	SignUpUser(ctx context.Context, user *User) (*User, error)
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	Logout(ctx context.Context, userID int64) error
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
}
