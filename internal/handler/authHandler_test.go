package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nicewook/gocore/internal/config"
	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
	"github.com/nicewook/gocore/pkg/security"
	"github.com/nicewook/gocore/pkg/validatorutil"
)

var (
	authConfig = &config.Config{
		Secure: config.SecureConfig{
			JWT: config.JWTConfig{
				Cookie: config.CookieConfig{
					Secure:   false,
					HTTPOnly: true,
					SameSite: "Lax",
					Domain:   "localhost",
				},
			},
		},
	}
)

func TestAuthHandler_SignUpUser(t *testing.T) {
	tests := []struct {
		name           string
		signupRequest  string
		mockInput      *domain.User
		mockReturn     interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			signupRequest:  `{"email":"john@example.com","password":"password123456"}`,
			mockInput:      &domain.User{Email: "john@example.com", Password: "password123456"},
			mockReturn:     &domain.User{Email: "john@example.com"},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "InvalidInput",
			signupRequest:  `{"email":"","password":""}`,
			mockInput:      nil,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "AlreadyExists",
			signupRequest:  `{"email":"john@example.com","password":"password123456"}`,
			mockInput:      &domain.User{Email: "john@example.com", Password: "password123456"},
			mockReturn:     nil,
			mockError:      domain.ErrAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = validatorutil.NewValidator()

			req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(tt.signupRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockUseCase := new(mocks.AuthUseCase)
			if tt.mockInput != nil && tt.name != "InvalidInput" {
				mockUseCase.On("SignUpUser", mock.Anything, tt.mockInput).Return(tt.mockReturn, tt.mockError).Maybe()
			}

			handler := NewAuthHandler(e, mockUseCase, authConfig)

			err := handler.SignUpUser(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		loginRequest   string
		mockEmail      string
		mockPassword   string
		mockReturn     *domain.LoginResponse
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			loginRequest:   `{"email":"john@example.com","password":"password"}`,
			mockEmail:      "john@example.com",
			mockPassword:   "password",
			mockReturn:     &domain.LoginResponse{ID: 1, Email: "john@example.com", AccessToken: "jwt-token", RefreshToken: "refresh-token", RefreshTokenExpiration: time.Now().Add(24 * time.Hour)},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Input",
			loginRequest:   `{"email":"","password":""}`,
			mockEmail:      "",
			mockPassword:   "",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Credentials",
			loginRequest:   `{"email":"john@example.com","password":"wrong"}`,
			mockEmail:      "john@example.com",
			mockPassword:   "wrong",
			mockReturn:     nil,
			mockError:      domain.ErrInvalidInput,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Internal Error",
			loginRequest:   `{"email":"john@example.com","password":"password"}`,
			mockEmail:      "john@example.com",
			mockPassword:   "password",
			mockReturn:     nil,
			mockError:      domain.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			e.Validator = validatorutil.NewValidator()

			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.loginRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockAuthUseCase := new(mocks.AuthUseCase)
			if tt.mockEmail != "" && tt.mockPassword != "" && tt.name != "Invalid Input" {
				mockAuthUseCase.On("Login", mock.Anything, tt.mockEmail, tt.mockPassword).Return(tt.mockReturn, tt.mockError)
			}

			h := NewAuthHandler(e, mockAuthUseCase, authConfig)

			// Act
			err := h.Login(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockAuthUseCase.AssertExpectations(t)

			// 성공 케이스가 아니면 쿠키 검증을 스킵
			if tt.expectedStatus != http.StatusOK {
				return
			}

			// 성공 케이스에서는 쿠키가 설정되었는지 확인
			cookies := rec.Result().Cookies()

			var refreshTokenCookie *http.Cookie
			for _, cookie := range cookies {
				if cookie.Name == refreshTokenCookieName {
					refreshTokenCookie = cookie
					break
				}
			}

			assert.NotNil(t, refreshTokenCookie, "refresh_token 쿠키가 설정되어야 합니다")
			assert.Equal(t, tt.mockReturn.RefreshToken, refreshTokenCookie.Value)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	// 테스트 데이터 정의
	tests := []struct {
		name           string
		setupAuth      func(c echo.Context)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupAuth: func(c echo.Context) {
				// JWT 토큰을 모킹하는 대신 user_id를 직접 설정
				c.Set("user", &jwt.Token{
					Claims: jwt.MapClaims{
						"user_id": float64(1),
						"email":   "john@example.com",
						"roles":   []interface{}{"User"},
					},
				})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Successfully logged out. Please remove the access token from your client storage.","status":"success"}`,
		},
		// 다른 테스트 케이스는 필요 시 추가
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			e.Validator = validatorutil.NewValidator()

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 인증 설정
			if tt.setupAuth != nil {
				tt.setupAuth(c)
			}

			mockUseCase := new(mocks.AuthUseCase)
			mockUseCase.On("Logout", mock.Anything, int64(1)).Return(nil)

			handler := NewAuthHandler(e, mockUseCase, authConfig)

			// Act
			err := handler.Logout(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			// 쿠키가 만료되었는지 확인
			cookies := rec.Result().Cookies()
			var refreshCookie *http.Cookie
			for _, cookie := range cookies {
				if cookie.Name == refreshTokenCookieName {
					refreshCookie = cookie
					break
				}
			}

			assert.NotNil(t, refreshCookie, "리프레시 토큰 쿠키가 응답에 있어야 합니다")
			assert.True(t, refreshCookie.Expires.Before(time.Now()), "쿠키가 만료되어야 합니다")

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		cookieValue    string
		mockReturn     *domain.LoginResponse
		mockError      error
		expectedStatus int
	}{
		{
			name:        "Success",
			cookieValue: "valid-refresh-token",
			mockReturn: &domain.LoginResponse{
				ID:                     1,
				Email:                  "john@example.com",
				AccessToken:            "new-access-token",
				RefreshToken:           "new-refresh-token",
				RefreshTokenExpiration: time.Now().Add(24 * time.Hour),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Token",
			cookieValue:    "invalid-token",
			mockReturn:     nil,
			mockError:      security.ErrInvalidToken,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "User Not Found",
			cookieValue:    "valid-token-unknown-user",
			mockReturn:     nil,
			mockError:      domain.ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Database Error",
			cookieValue:    "valid-token-db-error",
			mockReturn:     nil,
			mockError:      domain.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Echo
			e := echo.New()
			e.Validator = validatorutil.NewValidator()

			// 쿠키 설정
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh-token", nil)
			req.AddCookie(&http.Cookie{
				Name:  refreshTokenCookieName,
				Value: tt.cookieValue,
			})
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 모킹
			mockUseCase := new(mocks.AuthUseCase)
			mockUseCase.On("RefreshToken", mock.Anything, tt.cookieValue).Return(tt.mockReturn, tt.mockError)

			handler := NewAuthHandler(e, mockUseCase, authConfig)

			// 테스트 실행
			err := handler.RefreshToken(c)

			// 결과 검증
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// 성공 케이스는 추가 검증
			if tt.expectedStatus == http.StatusOK {
				// 응답 본문 확인
				var response domain.LoginResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturn.AccessToken, response.AccessToken)
				assert.Equal(t, tt.mockReturn.ID, response.ID)
				assert.Equal(t, tt.mockReturn.Email, response.Email)

				// 새 리프레시 토큰 쿠키 확인
				cookies := rec.Result().Cookies()
				var refreshCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == refreshTokenCookieName {
						refreshCookie = cookie
						break
					}
				}
				assert.NotNil(t, refreshCookie, "리프레시 토큰 쿠키가 응답에 있어야 합니다")
				assert.Equal(t, tt.mockReturn.RefreshToken, refreshCookie.Value)
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_CreateLogoutCookie(t *testing.T) {
	// Setup
	h := NewAuthHandler(echo.New(), nil, authConfig)

	// Act
	cookie := h.CreateLogoutCookie()

	// Assert
	assert.Equal(t, refreshTokenCookieName, cookie.Name)
	assert.Equal(t, "", cookie.Value)
	assert.True(t, cookie.Expires.Before(time.Now()), "쿠키가 만료되어야 합니다")
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, authConfig.Secure.JWT.Cookie.Domain, cookie.Domain)
	assert.Equal(t, authConfig.Secure.JWT.Cookie.Secure, cookie.Secure)
	assert.Equal(t, authConfig.Secure.JWT.Cookie.HTTPOnly, cookie.HttpOnly)
}

// TestSignUpRequestValidation 테스트는 SignUpRequest 구조체의 유효성 검사 태그를 확인합니다
func TestSignUpRequestValidation(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedCode  int
		expectedError bool
	}{
		{
			name:          "Valid Request",
			requestBody:   `{"email":"test@example.com","password":"password123"}`,
			expectedCode:  http.StatusCreated,
			expectedError: false,
		},
		{
			name:          "Empty Email",
			requestBody:   `{"email":"","password":"password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Invalid Email Format",
			requestBody:   `{"email":"invalid-email","password":"password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Empty Password",
			requestBody:   `{"email":"test@example.com","password":""}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Password Too Short",
			requestBody:   `{"email":"test@example.com","password":"short"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Missing Email Field",
			requestBody:   `{"password":"password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Missing Password Field",
			requestBody:   `{"email":"test@example.com"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Malformed JSON",
			requestBody:   `{"email":"test@example.com","password":}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Restricted Email Domain - Hotmail",
			requestBody:   `{"email":"test@hotmail.com","password":"password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Echo
			e := echo.New()
			e.Validator = validatorutil.NewValidator()

			// Setup request
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup mock
			mockUseCase := new(mocks.AuthUseCase)

			// 유효한 요청인 경우에만 모킹
			if !tt.expectedError {
				mockUser := &domain.User{
					Email:    "test@example.com",
					Password: "password123",
				}
				mockCreatedUser := &domain.User{
					ID:    1,
					Email: "test@example.com",
				}
				mockUseCase.On("SignUpUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.Email == mockUser.Email && u.Password == mockUser.Password
				})).Return(mockCreatedUser, nil)
			}

			// Create handler
			handler := NewAuthHandler(e, mockUseCase, authConfig)

			// Call handler method
			err := handler.SignUpUser(c)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)

			// 실패 케이스에서는 에러 응답이 있는지 확인
			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			} else {
				// 성공 케이스에서는 모킹 함수 호출 확인
				mockUseCase.AssertExpectations(t)

				// 응답 내용 확인
				var response domain.SignUpResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), response.ID)
				assert.Equal(t, "test@example.com", response.Email)
			}
		})
	}
}
