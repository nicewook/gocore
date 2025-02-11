package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		mockInput  *domain.User
		mockReturn *domain.User
		mockError  error
		expected   *domain.User
		expectErr  error
	}{
		{
			name:       "Success",
			mockInput:  &domain.User{Name: "John", Email: "john@example.com"},
			mockReturn: &domain.User{Name: "John", Email: "john@example.com"},
			mockError:  nil,
			expected:   &domain.User{Name: "John", Email: "john@example.com"},
			expectErr:  nil,
		},
		{
			name:      "InvalidInput",
			mockInput: &domain.User{Name: "", Email: ""},
			mockError: domain.ErrInvalidInput,
			expected:  nil,
			expectErr: domain.ErrInvalidInput,
		},
		{
			name:      "AlreadyExists",
			mockInput: &domain.User{Name: "John", Email: "john@example.com"},
			mockError: domain.ErrAlreadyExists,
			expected:  nil,
			expectErr: domain.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			mockRepo.On("Save", tt.mockInput).Return(tt.mockReturn, tt.mockError).Maybe()

			uc := NewUserUseCase(mockRepo)
			result, err := uc.CreateUser(tt.mockInput)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}

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
			name:      "User Not Found",
			inputID:   2,
			mockError: domain.ErrNotFound,
			expected:  nil,
			expectErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			mockRepo.On("GetByID", tt.inputID).Return(tt.mockReturn, tt.mockError)

			uc := NewUserUseCase(mockRepo)
			result, err := uc.GetByID(tt.inputID)

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
			mockRepo.On("GetAll").Return(tt.mockReturn, tt.mockError)

			uc := NewUserUseCase(mockRepo)
			result, err := uc.GetAll()

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
