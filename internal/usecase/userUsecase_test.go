package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
)

func TestGetByID(t *testing.T) {
	tests := []struct {
		name       string
		inputID    int64
		mockReturn *domain.User
		mockError  error
		expected   *domain.User
		expectErr  error
	}{
		{
			name:       "User Found",
			inputID:    1,
			mockReturn: &domain.User{ID: 1, Name: "John", Email: "john@example.com"},
			mockError:  nil,
			expected:   &domain.User{ID: 1, Name: "John", Email: "john@example.com"},
			expectErr:  nil,
		},
		{
			name:       "User Not Found",
			inputID:    2,
			mockReturn: nil,
			mockError:  domain.ErrNotFound,
			expected:   nil,
			expectErr:  domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			mockRepo.On("GetByID", mock.Anything, tt.inputID).Return(tt.mockReturn, tt.mockError)

			uc := NewUserUseCase(mockRepo)
			ctx := context.Background()
			result, err := uc.GetByID(ctx, tt.inputID)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserGetAll(t *testing.T) {
	repo := new(mocks.UserRepository)
	useCase := NewUserUseCase(repo)
	ctx := context.Background()

	t.Run("사용자 목록 조회 성공", func(t *testing.T) {
		// Mock data
		userList := []domain.User{
			{ID: int64(1), Name: "User 1", Email: "user1@example.com"},
			{ID: int64(2), Name: "User 2", Email: "user2@example.com"},
		}

		expectedResponse := &domain.GetAllResponse{
			Users:      userList,
			TotalCount: int64(len(userList)),
			Offset:     0,
			Limit:      10,
			HasMore:    false,
		}

		req := &domain.GetAllRequest{
			Limit: 10,
		}

		// Setup mock
		repo.On("GetAll", ctx, req).Return(expectedResponse, nil)

		// Call use case method
		result, err := useCase.GetAll(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		assert.Len(t, result.Users, 2)
		assert.Equal(t, int64(2), result.TotalCount)
		assert.Equal(t, false, result.HasMore)

		repo.AssertExpectations(t)
	})

	t.Run("필터링 기반 사용자 조회", func(t *testing.T) {
		// Mock data for filtered response
		filteredUser := []domain.User{
			{ID: int64(1), Name: "John Doe", Email: "john@example.com"},
		}

		expectedResponse := &domain.GetAllResponse{
			Users:      filteredUser,
			TotalCount: int64(len(filteredUser)),
			Offset:     0,
			Limit:      10,
			HasMore:    false,
		}

		// 이름 필터링 요청
		nameReq := &domain.GetAllRequest{
			Name:  "John",
			Limit: 10,
		}

		// Setup mock
		repo.On("GetAll", ctx, nameReq).Return(expectedResponse, nil)

		// Call use case method
		result, err := useCase.GetAll(ctx, nameReq)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		assert.Len(t, result.Users, 1)
		assert.Equal(t, "John Doe", result.Users[0].Name)

		repo.AssertExpectations(t)
	})

	t.Run("페이지네이션 테스트", func(t *testing.T) {
		// First page data
		userList := []domain.User{
			{ID: int64(1), Name: "User 1", Email: "user1@example.com"},
			{ID: int64(2), Name: "User 2", Email: "user2@example.com"},
		}

		expectedResponse := &domain.GetAllResponse{
			Users:      userList,
			TotalCount: 3, // 총 3명이지만 페이지당 2명씩 표시
			Offset:     0,
			Limit:      2,
			HasMore:    true,
		}

		req := &domain.GetAllRequest{
			Offset: 0,
			Limit:  2,
		}

		// Setup mock
		repo.On("GetAll", ctx, req).Return(expectedResponse, nil)

		// Call use case method
		result, err := useCase.GetAll(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		assert.Len(t, result.Users, 2)
		assert.Equal(t, int64(3), result.TotalCount)
		assert.True(t, result.HasMore)

		repo.AssertExpectations(t)
	})

	t.Run("사용자를 찾을 수 없는 경우", func(t *testing.T) {
		// Reset mock for this specific test case
		repo := new(mocks.UserRepository)
		useCase := NewUserUseCase(repo)

		req := &domain.GetAllRequest{
			Limit: 10,
		}

		// Setup mock with exact request matching
		repo.On("GetAll", ctx, req).Return(nil, domain.ErrNotFound)

		// Call use case method
		result, err := useCase.GetAll(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, result)

		repo.AssertExpectations(t)
	})
}
