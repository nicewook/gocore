package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/domain"
)

type OrderHandler struct {
	orderUseCase domain.OrderUseCase
}

func NewOrderHandler(e *echo.Echo, orderUseCase domain.OrderUseCase) *OrderHandler {
	handler := &OrderHandler{orderUseCase: orderUseCase}

	group := e.Group("/orders")
	group.POST("", handler.CreateOrder)
	group.GET("", handler.GetAll)
	group.GET("/:id", handler.GetByID)

	return handler
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var order domain.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	if order.UserID <= 0 || order.ProductID <= 0 || order.Quantity <= 0 || order.TotalPriceInKRW <= 0 {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	ctx := c.Request().Context()
	createdOrder, err := h.orderUseCase.CreateOrder(ctx, &order)
	if err == nil {
		return c.JSON(http.StatusCreated, createdOrder)
	}
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return c.JSON(http.StatusBadRequest, ErrResponse(err))
	case errors.Is(err, domain.ErrAlreadyExists):
		return c.JSON(http.StatusConflict, ErrResponse(err))
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
