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

func TestGetAll(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []domain.User
		mockError  error
		expected   []domain.User
		expectErr  error
	}{
		{
			name: "Users Found",
			mockReturn: []domain.User{
				{ID: 1, Name: "John", Email: "john@example.com"},
				{ID: 2, Name: "Jane", Email: "jane@example.com"},
			},
			mockError: nil,
			expected: []domain.User{
				{ID: 1, Name: "John", Email: "john@example.com"},
				{ID: 2, Name: "Jane", Email: "jane@example.com"},
			},
			expectErr: nil,
		},
		{
			name:      "No Users Found",
			mockError: domain.ErrNotFound,
			expected:  nil,
			expectErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			mockRepo.On("GetAll", mock.Anything).Return(tt.mockReturn, tt.mockError)

			uc := NewUserUseCase(mockRepo)
			ctx := context.Background()
			result, err := uc.GetAll(ctx)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
