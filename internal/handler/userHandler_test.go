package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
	"github.com/nicewook/gocore/pkg/validatorutil"
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

	// Setup validator for testing
	e := echo.New()
	e.Validator = validatorutil.NewValidator()

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
			mockError:      nil, // Error will be caught by validation
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidInput.Error()),
		},
		{
			name:           "Negative ID",
			pathParam:      "-1",
			mockReturn:     nil,
			mockError:      nil, // Error will be caught by validation
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidInput.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// echo context 생성
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.pathParam, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParam)

			// mock 생성, 설정 및 핸들러 생성
			mockUseCase := new(mocks.UserUseCase)

			// Only set up the mock expectation for valid IDs that will pass validation
			if tt.pathParam != "invalid" && tt.pathParam != "-1" {
				id, _ := strconv.ParseInt(tt.pathParam, 10, 64)
				mockUseCase.On("GetByID", mock.Anything, id).Return(tt.mockReturn, tt.mockError)
			}

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
	adminRoles := []string{domain.RoleAdmin}

	// Setup validator for testing
	e := echo.New()
	e.Validator = validatorutil.NewValidator()

	// Sample users for testing
	allUsers := []domain.User{
		{ID: 1, Name: "John", Email: "john@example.com", Roles: userRoles},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Roles: managerRoles},
		{ID: 3, Name: "Admin", Email: "admin@example.com", Roles: adminRoles},
	}

	// Filter for name "John"
	johnUsers := []domain.User{
		{ID: 1, Name: "John", Email: "john@example.com", Roles: userRoles},
	}

	// Filter for role "Manager"
	managerUsers := []domain.User{
		{ID: 2, Name: "Jane", Email: "jane@example.com", Roles: managerRoles},
	}

	// Filter for email containing "admin"
	adminUsers := []domain.User{
		{ID: 3, Name: "Admin", Email: "admin@example.com", Roles: adminRoles},
	}

	// Paginated users (first page, limit 2)
	paginatedUsers := []domain.User{
		{ID: 1, Name: "John", Email: "john@example.com", Roles: userRoles},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Roles: managerRoles},
	}

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockParams     *domain.GetAllRequest
		mockReturn     *domain.GetAllResponse
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Get all users successfully",
			queryParams: map[string]string{},
			mockParams:  &domain.GetAllRequest{Limit: 10},
			mockReturn: &domain.GetAllResponse{
				Users:      allUsers,
				TotalCount: int64(len(allUsers)),
				Offset:     0,
				Limit:      10,
				HasMore:    false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`{
				"users": [
					{"id":1,"name":"John","email":"john@example.com","roles":["%s"]},
					{"id":2,"name":"Jane","email":"jane@example.com","roles":["%s","%s"]},
					{"id":3,"name":"Admin","email":"admin@example.com","roles":["%s"]}
				],
				"total_count": 3,
				"offset": 0,
				"limit": 10,
				"has_more": false
			}`, domain.RoleUser, domain.RoleManager, domain.RoleUser, domain.RoleAdmin),
		},
		{
			name: "Filter by name",
			queryParams: map[string]string{
				"name": "John",
			},
			mockParams: &domain.GetAllRequest{
				Name:  "John",
				Limit: 10,
			},
			mockReturn: &domain.GetAllResponse{
				Users:      johnUsers,
				TotalCount: int64(len(johnUsers)),
				Offset:     0,
				Limit:      10,
				HasMore:    false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`{
				"users": [
					{"id":1,"name":"John","email":"john@example.com","roles":["%s"]}
				],
				"total_count": 1,
				"offset": 0,
				"limit": 10,
				"has_more": false
			}`, domain.RoleUser),
		},
		{
			name: "Filter by email",
			queryParams: map[string]string{
				"email": "admin",
			},
			mockParams: &domain.GetAllRequest{
				Email: "admin",
				Limit: 10,
			},
			mockReturn: &domain.GetAllResponse{
				Users:      adminUsers,
				TotalCount: int64(len(adminUsers)),
				Offset:     0,
				Limit:      10,
				HasMore:    false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`{
				"users": [
					{"id":3,"name":"Admin","email":"admin@example.com","roles":["%s"]}
				],
				"total_count": 1,
				"offset": 0,
				"limit": 10,
				"has_more": false
			}`, domain.RoleAdmin),
		},
		{
			name: "Filter by role",
			queryParams: map[string]string{
				"roles": "Manager",
			},
			mockParams: &domain.GetAllRequest{
				Roles: "Manager",
				Limit: 10,
			},
			mockReturn: &domain.GetAllResponse{
				Users:      managerUsers,
				TotalCount: int64(len(managerUsers)),
				Offset:     0,
				Limit:      10,
				HasMore:    false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`{
				"users": [
					{"id":2,"name":"Jane","email":"jane@example.com","roles":["%s","%s"]}
				],
				"total_count": 1,
				"offset": 0,
				"limit": 10,
				"has_more": false
			}`, domain.RoleManager, domain.RoleUser),
		},
		{
			name: "Pagination",
			queryParams: map[string]string{
				"offset": "0",
				"limit":  "2",
			},
			mockParams: &domain.GetAllRequest{
				Offset: 0,
				Limit:  2,
			},
			mockReturn: &domain.GetAllResponse{
				Users:      paginatedUsers,
				TotalCount: int64(len(allUsers)),
				Offset:     0,
				Limit:      2,
				HasMore:    true,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`{
				"users": [
					{"id":1,"name":"John","email":"john@example.com","roles":["%s"]},
					{"id":2,"name":"Jane","email":"jane@example.com","roles":["%s","%s"]}
				],
				"total_count": 3,
				"offset": 0,
				"limit": 2,
				"has_more": true
			}`, domain.RoleUser, domain.RoleManager, domain.RoleUser),
		},
		{
			name: "Pagination - Second Page",
			queryParams: map[string]string{
				"offset": "2", // 3번째 항목부터 (0-indexed)
				"limit":  "2",
			},
			mockParams: &domain.GetAllRequest{
				Offset: 2,
				Limit:  2,
			},
			mockReturn: &domain.GetAllResponse{
				Users:      []domain.User{{ID: 3, Name: "Admin", Email: "admin@example.com", Roles: adminRoles}},
				TotalCount: int64(len(allUsers)),
				Offset:     2,
				Limit:      2,
				HasMore:    false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`{
				"users": [
					{"id":3,"name":"Admin","email":"admin@example.com","roles":["%s"]}
				],
				"total_count": 3,
				"offset": 2,
				"limit": 2,
				"has_more": false
			}`, domain.RoleAdmin),
		},
		{
			name: "Invalid pagination parameters",
			queryParams: map[string]string{
				"offset": "-1", // Negative offset
			},
			mockParams:     nil, // No mock call expected due to validation error
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidInput.Error()),
		},
		{
			name: "No users found",
			queryParams: map[string]string{
				"name": "NonExistentUser",
			},
			mockParams: &domain.GetAllRequest{
				Name:  "NonExistentUser",
				Limit: 10,
			},
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusOK, // 빈 목록과 페이지네이션 정보를 반환
			expectedBody: `{
				"users": [],
				"total_count": 0,
				"offset": 0,
				"limit": 10,
				"has_more": false
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create query string
			queryString := ""
			if len(tt.queryParams) > 0 {
				queryValues := make([]string, 0)
				for key, value := range tt.queryParams {
					queryValues = append(queryValues, key+"="+value)
				}
				queryString = "?" + strings.Join(queryValues, "&")
			}

			// Create Echo context
			req := httptest.NewRequest(http.MethodGet, "/users"+queryString, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Create mock UseCase
			mockUseCase := new(mocks.UserUseCase)

			// 유효성 검사 오류가 예상되는 경우 mock 설정을 하지 않음
			if tt.mockParams != nil {
				// Use mock.MatchedBy to verify the request parameters
				mockUseCase.On("GetAll", mock.Anything, mock.MatchedBy(func(req *domain.GetAllRequest) bool {
					reqParam := tt.mockParams
					return req.Name == reqParam.Name &&
						req.Email == reqParam.Email &&
						req.Roles == reqParam.Roles &&
						req.Offset == reqParam.Offset &&
						req.Limit == reqParam.Limit
				})).Return(tt.mockReturn, tt.mockError)
			}

			// Create handler with mock UseCase
			handler := NewUserHandler(e, mockUseCase)

			// Call the handler
			err := handler.GetAll(c)
			assert.NoError(t, err)

			// Check the response
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			// Verify mock expectations
			mockUseCase.AssertExpectations(t)
		})
	}
}
