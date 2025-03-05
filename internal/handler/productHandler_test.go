package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/domain/mocks"
)

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		mockInput      *domain.Product
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			input:          `{"name":"Product1","price_in_krw":100}`,
			mockInput:      &domain.Product{Name: "Product1", PriceInKRW: 100},
			mockReturn:     &domain.Product{ID: 1, Name: "Product1", PriceInKRW: 100},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"name":"Product1","price_in_krw":100}`,
		},
		{
			name:           "InvalidInput",
			input:          `{"name":"","price_in_krw":0}`,
			mockInput:      &domain.Product{Name: "", PriceInKRW: 0},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid input"}`,
		},
		{
			name:           "InvalidJSON",
			input:          `{"name":"Product1",price_in_krw:100}`,
			mockInput:      nil,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid input"}`,
		},
		{
			name:           "AlreadyExists",
			input:          `{"name":"Product1","price_in_krw":100}`,
			mockInput:      &domain.Product{Name: "Product1", PriceInKRW: 100},
			mockReturn:     nil,
			mockError:      domain.ErrAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"already exists"}`,
		},
		{
			name:           "InternalError",
			input:          `{"name":"Product1","price_in_krw":100}`,
			mockInput:      &domain.Product{Name: "Product1", PriceInKRW: 100},
			mockReturn:     nil,
			mockError:      domain.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}`,
		},
		{
			name:           "InvalidInputFromUseCase",
			input:          `{"name":"Product1","price_in_krw":100}`,
			mockInput:      &domain.Product{Name: "Product1", PriceInKRW: 100},
			mockReturn:     nil,
			mockError:      domain.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid input"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(tt.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			mockUseCase := new(mocks.ProductUseCase)
			if tt.mockInput != nil {
				mockUseCase.On("CreateProduct", mock.Anything, tt.mockInput).Return(tt.mockReturn, tt.mockError).Maybe()
			}
			handler := NewProductHandler(e, mockUseCase)
			err := handler.CreateProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// 응답 본문 검증
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestGetProductByID(t *testing.T) {
	tests := []struct {
		name           string
		pathParam      string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get product by ID successfully",
			pathParam:      "1",
			mockReturn:     &domain.Product{ID: 1, Name: "Product1", PriceInKRW: 100},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Product1","price_in_krw":100}`,
		},
		{
			name:           "Product not found",
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
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.pathParam, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParam)
			mockUseCase := new(mocks.ProductUseCase)
			mockUseCase.On("GetByID", mock.Anything, mock.Anything).Return(tt.mockReturn, tt.mockError).Maybe()
			handler := NewProductHandler(e, mockUseCase)
			err := handler.GetByID(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	tests := []struct {
		name           string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Get all products successfully",
			mockReturn: []domain.Product{
				{ID: 1, Name: "Product1", PriceInKRW: 100},
				{ID: 2, Name: "Product2", PriceInKRW: 200},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"Product1","price_in_krw":100},{"id":2,"name":"Product2","price_in_krw":200}]`,
		},
		{
			name:           "No products found",
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			mockUseCase := new(mocks.ProductUseCase)
			mockUseCase.On("GetAll", mock.Anything).Return(tt.mockReturn, tt.mockError)
			handler := NewProductHandler(e, mockUseCase)
			err := handler.GetAll(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockUseCase.AssertExpectations(t)
		})
	}
}
