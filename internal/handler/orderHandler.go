package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/middlewares"
)

type OrderHandler struct {
	orderUseCase domain.OrderUseCase
}

func NewOrderHandler(e *echo.Echo, orderUseCase domain.OrderUseCase) *OrderHandler {
	handler := &OrderHandler{orderUseCase: orderUseCase}

	group := e.Group("/orders")
	group.POST("", handler.CreateOrder, middlewares.AllowRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser))
	group.GET("", handler.GetAll, middlewares.AllowRoles(domain.RoleAdmin, domain.RoleManager))
	group.GET("/:id", handler.GetByID, middlewares.AllowRoles(domain.RoleAdmin, domain.RoleManager, domain.RoleUser))

	return handler
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	// 주문 데이터 바인딩
	order := new(domain.Order)
	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	// 입력값 검증
	if err := c.Validate(order); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	ctx := c.Request().Context()

	// 주문 생성 처리
	createdOrder, err := h.orderUseCase.CreateOrder(ctx, order)
	if err == nil {
		return c.JSON(http.StatusCreated, createdOrder)
	}

	// 에러 처리
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return c.JSON(http.StatusBadRequest, ErrResponse(err))
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(err))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}

func (h *OrderHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	ctx := c.Request().Context()
	order, err := h.orderUseCase.GetByID(ctx, int64(id))
	if err == nil {
		return c.JSON(http.StatusOK, order)
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(domain.ErrNotFound))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}

func (h *OrderHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	orders, err := h.orderUseCase.GetAll(ctx)
	if err == nil {
		return c.JSON(http.StatusOK, orders)
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(domain.ErrNotFound))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}
