package security

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrParsingKey   = errors.New("error parsing RSA key")
)

// TokenType defines the type of token
type TokenType string

const (
	// AccessToken is used for API access
	AccessToken TokenType = "access"
	// RefreshToken is used to get new access tokens
	RefreshToken TokenType = "refresh"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID int64     `json:"user_id"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// ParseRSAPrivateKeyFromPEM parses a PEM encoded RSA private key
func ParseRSAPrivateKeyFromPEM(key string) (*rsa.PrivateKey, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		return nil, ErrParsingKey
	}
	return privateKey, nil
}

// ParseRSAPublicKeyFromPEM parses a PEM encoded RSA public key
func ParseRSAPublicKeyFromPEM(key string) (*rsa.PublicKey, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		return nil, ErrParsingKey
	}
	return publicKey, nil
}

// GenerateToken creates a new JWT token for a user using RSA private key
func GenerateToken(userID int64, email string, roles []string, privateKey *rsa.PrivateKey, expirationTime time.Duration, tokenType TokenType) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateAccessToken creates a new access token
func GenerateAccessToken(userID int64, email string, roles []string, privateKey *rsa.PrivateKey, expirationTime time.Duration) (string, error) {
	return GenerateToken(userID, email, roles, privateKey, expirationTime, AccessToken)
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken(userID int64, email string, roles []string, privateKey *rsa.PrivateKey, expirationTime time.Duration) (string, error) {
	return GenerateToken(userID, email, roles, privateKey, expirationTime, RefreshToken)
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string, publicKey *rsa.PublicKey) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 서명 방식 확인
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func ValidateAccessToken(tokenString string, publicKey *rsa.PublicKey) (*JWTClaims, error) {
	claims, err := ValidateToken(tokenString, publicKey)
	if err != nil {
		return nil, err
	}

	if claims.Type != AccessToken {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func ValidateRefreshToken(tokenString string, publicKey *rsa.PublicKey) (*JWTClaims, error) {
	claims, err := ValidateToken(tokenString, publicKey)
	if err != nil {
		return nil, err
	}

	if claims.Type != RefreshToken {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
