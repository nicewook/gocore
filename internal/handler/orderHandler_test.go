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

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		mockInput      *domain.Order
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			input:          `{"user_id":1,"product_id":1,"quantity":2,"total_price_in_krw":2000}`,
			mockInput:      &domain.Order{UserID: 1, ProductID: 1, Quantity: 2, TotalPriceInKRW: 2000},
			mockReturn:     &domain.Order{ID: 1, UserID: 1, ProductID: 1, Quantity: 2, TotalPriceInKRW: 2000},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"user_id":1,"product_id":1,"quantity":2,"total_price_in_krw":2000,"created_at":""}`,
		},
		{
			name:           "InvalidInput",
			input:          `{"user_id":0,"product_id":0,"quantity":0,"total_price_in_krw":0}`,
			mockInput:      &domain.Order{UserID: 0, ProductID: 0, Quantity: 0, TotalPriceInKRW: 0},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid input"}`,
		},
		{
			name:           "InvalidJSON",
			input:          `{"user_id":1,"product_id":1,quantity:2,"total_price_in_krw":2000}`,
			mockInput:      nil,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid input"}`,
		},
		{
			name:           "AlreadyExists",
			input:          `{"user_id":1,"product_id":1,"quantity":2,"total_price_in_krw":2000}`,
			mockInput:      &domain.Order{UserID: 1, ProductID: 1, Quantity: 2, TotalPriceInKRW: 2000},
			mockReturn:     nil,
			mockError:      domain.ErrAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"already exists"}`,
		},
		{
			name:           "InternalError",
			input:          `{"user_id":1,"product_id":1,"quantity":2,"total_price_in_krw":2000}`,
			mockInput:      &domain.Order{UserID: 1, ProductID: 1, Quantity: 2, TotalPriceInKRW: 2000},
			mockReturn:     nil,
			mockError:      domain.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}`,
		},
		{
			name:           "InvalidInputFromUseCase",
			input:          `{"user_id":1,"product_id":1,"quantity":2,"total_price_in_krw":2000}`,
			mockInput:      &domain.Order{UserID: 1, ProductID: 1, Quantity: 2, TotalPriceInKRW: 2000},
			mockReturn:     nil,
			mockError:      domain.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid input"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(tt.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockUseCase := new(mocks.OrderUseCase)
			if tt.mockInput != nil {
				mockUseCase.On("CreateOrder", mock.Anything, tt.mockInput).Return(tt.mockReturn, tt.mockError).Maybe()
			}
			handler := NewOrderHandler(e, mockUseCase)

			err := handler.CreateOrder(c)
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

func TestGetOrderByID(t *testing.T) {
	tests := []struct {
		name           string
		pathParam      string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get order by ID successfully",
			pathParam:      "1",
			mockReturn:     &domain.Order{ID: 1, UserID: 1, ProductID: 1, Quantity: 2, TotalPriceInKRW: 2000},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"user_id":1,"product_id":1,"quantity":2,"total_price_in_krw":2000,"created_at":""}`,
		},
		{
			name:           "Order not found",
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
			req := httptest.NewRequest(http.MethodGet, "/orders/"+tt.pathParam, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParam)

			mockUseCase := new(mocks.OrderUseCase)
			mockUseCase.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(tt.mockReturn, tt.mockError).Maybe()
			handler := NewOrderHandler(e, mockUseCase)

			err := handler.GetByID(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestGetAllOrders(t *testing.T) {
	tests := []struct {
		name           string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Get all orders successfully",
			mockReturn: []domain.Order{
				{ID: 1, UserID: 1, ProductID: 1, Quantity: 2, TotalPriceInKRW: 2000},
				{ID: 2, UserID: 2, ProductID: 2, Quantity: 1, TotalPriceInKRW: 1000},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,

			expectedBody: `[
				{"id":1,"user_id":1,"product_id":1,"quantity":2,"total_price_in_krw":2000,"created_at":""},
				{"id":2,"user_id":2,"product_id":2,"quantity":1,"total_price_in_krw":1000,"created_at":""}
			]`,
		},
		{
			name:           "No orders found",
			mockReturn:     nil,
			mockError:      domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   fmt.Sprintf(`{"error":"%s"}`, domain.ErrNotFound.Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			mockUseCase := new(mocks.OrderUseCase)
			mockUseCase.On("GetAll", mock.Anything).Return(tt.mockReturn, tt.mockError)
			handler := NewOrderHandler(e, mockUseCase)
			err := handler.GetAll(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockUseCase.AssertExpectations(t)
		})
	}
}
