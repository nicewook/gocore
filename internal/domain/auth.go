package domain

import (
	"context"
	"time"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
