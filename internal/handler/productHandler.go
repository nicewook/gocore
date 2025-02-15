package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/domain"
)

type ProductHandler struct {
	productUseCase domain.ProductUseCase
}

func NewProductHandler(productUseCase domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	product := new(domain.Product)
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}
	if product.Name == "" || product.PriceInKRW <= 0 {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}
	_, err := h.productUseCase.CreateProduct(product)
	if err == nil {
		return c.JSON(http.StatusCreated, product)
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
	product, err := h.productUseCase.GetByID(int64(id))
	if err == nil {
		return c.JSON(http.StatusOK, product)
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(err))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}

func (h *ProductHandler) GetAll(c echo.Context) error {
	products, err := h.productUseCase.GetAll()
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
