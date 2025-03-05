package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name       string
		mockInput  *domain.Order
		mockReturn *domain.Order
		mockError  error
		expected   *domain.Order
		expectErr  error
	}{
		{
			name:       "Success",
			mockInput:  &domain.Order{UserID: 1, ProductID: 1, Quantity: 1, TotalPriceInKRW: 1000},
			mockReturn: &domain.Order{UserID: 1, ProductID: 1, Quantity: 1, TotalPriceInKRW: 1000},
			mockError:  nil,
			expected:   &domain.Order{UserID: 1, ProductID: 1, Quantity: 1, TotalPriceInKRW: 1000},
			expectErr:  nil,
		},
		{
			name:       "InvalidInput",
			mockInput:  &domain.Order{UserID: 0, ProductID: 0, Quantity: 0, TotalPriceInKRW: 0},
			mockReturn: nil,
			mockError:  domain.ErrInvalidInput,
			expected:   nil,
			expectErr:  domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.OrderRepository)
			mockRepo.On("Save", mock.Anything, tt.mockInput).Return(tt.mockReturn, tt.mockError).Maybe()

			uc := NewOrderUseCase(mockRepo)
			ctx := context.Background()
			result, err := uc.CreateOrder(ctx, tt.mockInput)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	tests := []struct {
		name       string
		inputID    int64
		mockReturn *domain.Order
		mockError  error
		expected   *domain.Order
		expectErr  error
	}{
		{
			name:       "Order Found",
			inputID:    1,
			mockReturn: &domain.Order{ID: 1, UserID: 1, ProductID: 1, Quantity: 1, TotalPriceInKRW: 1000},
			mockError:  nil,
			expected:   &domain.Order{ID: 1, UserID: 1, ProductID: 1, Quantity: 1, TotalPriceInKRW: 1000},
			expectErr:  nil,
		},
		{
			name:      "Order Not Found",
			inputID:   2,
			mockError: domain.ErrNotFound,
			expected:  nil,
			expectErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.OrderRepository)
			mockRepo.On("GetByID", mock.Anything, tt.inputID).Return(tt.mockReturn, tt.mockError)

			uc := NewOrderUseCase(mockRepo)
			ctx := context.Background()
			result, err := uc.GetByID(ctx, tt.inputID)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllOrders(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []domain.Order
		mockError  error
		expected   []domain.Order
		expectErr  error
	}{
		{
			name: "Orders Found",
			mockReturn: []domain.Order{
				{ID: 1, UserID: 1, ProductID: 1, Quantity: 1, TotalPriceInKRW: 1000},
				{ID: 2, UserID: 2, ProductID: 2, Quantity: 2, TotalPriceInKRW: 2000},
			},
			mockError: nil,
			expected: []domain.Order{
				{ID: 1, UserID: 1, ProductID: 1, Quantity: 1, TotalPriceInKRW: 1000},
				{ID: 2, UserID: 2, ProductID: 2, Quantity: 2, TotalPriceInKRW: 2000},
			},
			expectErr: nil,
		},
		{
			name:      "No Orders Found",
			mockError: domain.ErrNotFound,
			expected:  nil,
			expectErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.OrderRepository)
			mockRepo.On("GetAll", mock.Anything).Return(tt.mockReturn, tt.mockError)

			uc := NewOrderUseCase(mockRepo)
			ctx := context.Background()
			result, err := uc.GetAll(ctx)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
