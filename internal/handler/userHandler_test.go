package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
)

func TestErrResponse(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "InvalidInput",
			args: args{err: domain.ErrInvalidInput},
			want: map[string]string{"error": domain.ErrInvalidInput.Error()},
		},
		{
			name: "AlreadyExists",
			args: args{err: domain.ErrAlreadyExists},
			want: map[string]string{"error": domain.ErrAlreadyExists.Error()},
		},
		{
			name: "Internal",
			args: args{err: domain.ErrInternal},
			want: map[string]string{"error": domain.ErrInternal.Error()},
		},
		{
			name: "NilError",
			args: args{err: nil},
			want: map[string]string{"error": "unknown error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrResponse(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	// 사용자 역할 정의
	userRoles := []string{domain.RoleUser}
	adminRoles := []string{domain.RoleAdmin, domain.RoleUser}

	tests := []struct {
		name           string
		pathParam      string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get user by ID successfully",
			pathParam:      "1",
			mockReturn:     &domain.User{ID: 1, Name: "John", Email: "john@example.com", Roles: userRoles},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`{"id":1,"name":"John","email":"john@example.com","roles":["%s"]}`, domain.RoleUser),
		},
		{
			name:           "Get admin by ID successfully",
			pathParam:      "2",
			mockReturn:     &domain.User{ID: 2, Name: "Admin", Email: "admin@example.com", Roles: adminRoles},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`{"id":2,"name":"Admin","email":"admin@example.com","roles":["%s","%s"]}`, domain.RoleAdmin, domain.RoleUser),
		},
		{
			name:           "User not found",
			pathParam:      "1",
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrNotFound.Error()),
		},
		{
			name:           "Invalid ID format",
			pathParam:      "invalid",
			mockReturn:     nil,
			mockError:      domain.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidInput.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// echo context 생성
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.pathParam, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParam)

			// mcok 생성, 설정 및 핸들러 생성
			mockUseCase := new(mocks.UserUseCase)
			mockUseCase.On("GetByID", mock.Anything, mock.Anything).Return(tt.mockReturn, tt.mockError).Maybe()
			handler := NewUserHandler(e, mockUseCase)

			// 핸들러 실행 및 검증
			err := handler.GetByID(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			// mock 호출 검증
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestGetAll(t *testing.T) {
	// 사용자 역할 정의
	userRoles := []string{domain.RoleUser}
	managerRoles := []string{domain.RoleManager, domain.RoleUser}

	tests := []struct {
		name           string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Get all users successfully",
			mockReturn: []domain.User{
				{ID: 1, Name: "John", Email: "john@example.com", Roles: userRoles},
				{ID: 2, Name: "Jane", Email: "jane@example.com", Roles: managerRoles},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`[{"id":1,"name":"John","email":"john@example.com","roles":["%s"]},{"id":2,"name":"Jane","email":"jane@example.com","roles":["%s","%s"]}]`, domain.RoleUser, domain.RoleManager, domain.RoleUser),
		},
		{
			name:           "No users found",
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// echo context 생성
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// mock 생성, 설정 및 핸들러 생성
			mockUseCase := new(mocks.UserUseCase)
			mockUseCase.On("GetAll", mock.Anything).Return(tt.mockReturn, tt.mockError)
			handler := NewUserHandler(e, mockUseCase)

			// 핸들러 실행 및 검증
			err := handler.GetAll(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			// mock 호출 검증
			mockUseCase.AssertExpectations(t)
		})
	}
}
