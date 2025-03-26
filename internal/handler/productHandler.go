package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/middlewares"
)

type ProductHandler struct {
	productUseCase domain.ProductUseCase
}

func NewProductHandler(e *echo.Echo, productUseCase domain.ProductUseCase) *ProductHandler {
	handler := &ProductHandler{productUseCase: productUseCase}

	group := e.Group("/products")
	group.POST("", handler.CreateProduct, middlewares.AllowRoles(domain.RoleAdmin, domain.RoleManager))
	group.GET("", handler.GetAll, middlewares.AllowRoles(domain.RolePublic))
	group.GET("/:id", handler.GetByID, middlewares.AllowRoles(domain.RolePublic))

	return handler
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	product := new(domain.Product)
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}
	if product.Name == "" || product.PriceInKRW <= 0 {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	ctx := c.Request().Context()
	createdProduct, err := h.productUseCase.CreateProduct(ctx, product)
	if err == nil {
		return c.JSON(http.StatusCreated, createdProduct)
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

func (h *ProductHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	ctx := c.Request().Context()
	product, err := h.productUseCase.GetByID(ctx, int64(id))
	if err == nil {
		return c.JSON(http.StatusOK, product)
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(domain.ErrNotFound))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}

func (h *ProductHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	products, err := h.productUseCase.GetAll(ctx)
	if err == nil {
		return c.JSON(http.StatusOK, products)
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(err))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}
