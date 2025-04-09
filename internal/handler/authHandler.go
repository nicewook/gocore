package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/config"
	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/middlewares"
	"github.com/nicewook/gocore/pkg/contextutil"
	"github.com/nicewook/gocore/pkg/security"
)

const refreshTokenCookieName = "refresh_token"

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authUseCase domain.AuthUseCase
	config      *config.Config
}

func NewAuthHandler(e *echo.Echo, authUseCase domain.AuthUseCase, config *config.Config) *AuthHandler {
	handler := &AuthHandler{
		authUseCase: authUseCase,
		config:      config,
	}

	group := e.Group("/auth", middlewares.AllowRoles(domain.RolePublic))
	group.POST("/signup", handler.SignUpUser)
	group.POST("/login", handler.Login)
	group.POST("/refresh-token", handler.RefreshToken)
	group.POST("/logout", handler.Logout, middlewares.AllowRoles(
		domain.RoleAdmin, domain.RoleManager, domain.RoleUser))

	return handler
}

func (h *AuthHandler) SignUpUser(c echo.Context) error {
	req := new(domain.SignUpRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	if err := req.Validate(c); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request().Context()
	createdUser, err := h.authUseCase.SignUpUser(ctx, user)
	if err == nil {
		return c.JSON(http.StatusCreated, domain.SignUpResponse{
			ID:    createdUser.ID,
			Email: createdUser.Email,
		})
	}

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return c.JSON(http.StatusBadRequest, ErrResponse(err))
	case errors.Is(err, domain.ErrAlreadyExists):
		return c.JSON(http.StatusConflict, ErrResponse(err))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(domain.LoginRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	ctx := c.Request().Context()
	loginResponse, err := h.authUseCase.Login(ctx, req.Email, req.Password)
	if err == nil {
		// Set refresh token as HTTP-only cookie
		cookie := h.createRefreshTokenCookie(
			loginResponse.RefreshToken,
			loginResponse.RefreshTokenExpiration)
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, loginResponse)
	}

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return c.JSON(http.StatusUnauthorized, ErrResponse(errors.New("invalid email or password")))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}

// createRefreshTokenCookie creates a secure HTTP-only cookie for the refresh token
func (h *AuthHandler) createRefreshTokenCookie(tokenValue string, expiration time.Time) *http.Cookie {

	cookieConfig := h.config.Secure.JWT.Cookie

	// Parse SameSite value
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(cookieConfig.SameSite) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "lax":
		sameSite = http.SameSiteLaxMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    tokenValue,
		Path:     "/",
		Domain:   cookieConfig.Domain,
		Expires:  expiration,
		Secure:   cookieConfig.Secure,
		HttpOnly: cookieConfig.HTTPOnly,
		SameSite: sameSite,
	}

	return cookie
}

// CreateLogoutCookie creates a cookie that expires the refresh token
func (h *AuthHandler) CreateLogoutCookie() *http.Cookie {
	// 빈 값과 과거 만료일로 쿠키를 생성하여 쿠키를 삭제하는 효과를 냅니다
	return h.createRefreshTokenCookie("", time.Now().Add(-1*time.Hour))
}

// Logout handles user logout by invalidating tokens
// TODO: 향후 보안 강화를 위한 개선 사항
// 1. 리프레시 토큰 블랙리스트 구현:
//   - Redis 또는 다른 인메모리 저장소를 사용하여 무효화된 토큰 관리
//   - 토큰의 jti(JWT ID)를 블랙리스트에 추가하고 만료 시간까지 유지
//   - 모든 인증 요청에서 토큰이 블랙리스트에 있는지 확인
//
// 2. 다중 디바이스 지원:
//   - 사용자별로 발급된 모든 토큰 정보 저장
//   - 각 로그인 세션마다 고유 디바이스/세션 ID 발급
//   - JWT 클레임에 디바이스/세션 ID 포함
//   - 특정 디바이스만 로그아웃하거나 모든 디바이스에서 로그아웃 기능 구현
//
// 3. 보안 강화:
//   - 액세스 토큰 수명 단축 (5-15분)
//   - 중요 작업 수행 시 재인증 요구
//   - 비정상적인 접근 패턴 감지 및 차단
func (h *AuthHandler) Logout(c echo.Context) error {
	// 컨텍스트에서 사용자 ID 가져오기
	userID, _, _, err := contextutil.TokenToUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrResponse(domain.ErrUnauthorized))
	}

	ctx := c.Request().Context()
	// 로그아웃 처리 - 현재는 간단하게 구현
	err = h.authUseCase.Logout(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}

	// 리프레시 토큰 쿠키 만료시키기
	c.SetCookie(h.CreateLogoutCookie())

	// 클라이언트에 응답 - 프론트엔드에서 액세스 토큰을 삭제해야 함을 알림
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Successfully logged out. Please remove the access token from your client storage.",
		"status":  "success",
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// Get refresh token from cookie
	cookie, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrResponse(errors.New("refresh token not found")))
	}

	ctx := c.Request().Context()
	// Call use case to refresh tokens
	loginResponse, err := h.authUseCase.RefreshToken(ctx, cookie.Value)
	if err != nil {
		switch {
		case errors.Is(err, security.ErrInvalidToken), errors.Is(err, security.ErrExpiredToken):
			return c.JSON(http.StatusUnauthorized, ErrResponse(errors.New("invalid or expired refresh token")))
		case errors.Is(err, domain.ErrUnauthorized):
			return c.JSON(http.StatusUnauthorized, ErrResponse(domain.ErrUnauthorized))
		default:
			return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
		}
	}

	// Set new refresh token as HTTP-only cookie
	newCookie := h.createRefreshTokenCookie(loginResponse.RefreshToken, loginResponse.RefreshTokenExpiration)
	c.SetCookie(newCookie)

	// Return new access token
	return c.JSON(http.StatusOK, loginResponse)
}
