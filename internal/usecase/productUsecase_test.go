package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
)

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name       string
		mockInput  *domain.Product
		mockReturn *domain.Product
		mockError  error
		expected   *domain.Product
		expectErr  error
	}{
		{
			name:       "Success",
			mockInput:  &domain.Product{Name: "Product1", PriceInKRW: 100},
			mockReturn: &domain.Product{Name: "Product1", PriceInKRW: 100},
			mockError:  nil,
			expected:   &domain.Product{Name: "Product1", PriceInKRW: 100},
			expectErr:  nil,
		},
		{
			name:      "InvalidInput",
			mockInput: &domain.Product{Name: "", PriceInKRW: 0},
			mockError: domain.ErrInvalidInput,
			expected:  nil,
			expectErr: domain.ErrInvalidInput,
		},
		{
			name:      "AlreadyExists",
			mockInput: &domain.Product{Name: "Product1", PriceInKRW: 100},
			mockError: domain.ErrAlreadyExists,
			expected:  nil,
			expectErr: domain.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.ProductRepository)
			mockRepo.On("Save", tt.mockInput).Return(tt.mockReturn, tt.mockError).Maybe()
			uc := NewProductUseCase(mockRepo)
			result, err := uc.CreateProduct(tt.mockInput)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetProductByID(t *testing.T) {
	tests := []struct {
		name       string
		inputID    int64
		mockReturn *domain.Product
		mockError  error
		expected   *domain.Product
		expectErr  error
	}{
		{
			name:       "Product Found",
			inputID:    1,
			mockReturn: &domain.Product{ID: 1, Name: "Product1", PriceInKRW: 100},
			mockError:  nil,
			expected:   &domain.Product{ID: 1, Name: "Product1", PriceInKRW: 100},
			expectErr:  nil,
		},
		{
			name:      "Product Not Found",
			inputID:   2,
			mockError: domain.ErrNotFound,
			expected:  nil,
			expectErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.ProductRepository)
			mockRepo.On("GetByID", tt.inputID).Return(tt.mockReturn, tt.mockError)
			uc := NewProductUseCase(mockRepo)
			result, err := uc.GetByID(tt.inputID)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []domain.Product
		mockError  error
		expected   []domain.Product
		expectErr  error
	}{
		{
			name: "Products Found",
			mockReturn: []domain.Product{
				{ID: 1, Name: "Product1", PriceInKRW: 100},
				{ID: 2, Name: "Product2", PriceInKRW: 200},
			},
			mockError: nil,
			expected: []domain.Product{
				{ID: 1, Name: "Product1", PriceInKRW: 100},
				{ID: 2, Name: "Product2", PriceInKRW: 200},
			},
			expectErr: nil,
		},
		{
			name:      "No Products Found",
			mockError: domain.ErrNotFound,
			expected:  nil,
			expectErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.ProductRepository)
			mockRepo.On("GetAll").Return(tt.mockReturn, tt.mockError)
			uc := NewProductUseCase(mockRepo)
			result, err := uc.GetAll()
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
