package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/domain"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(e *echo.Echo, userUseCase domain.UserUseCase) *UserHandler {
	handler := &UserHandler{userUseCase: userUseCase}

	group := e.Group("/users")
	group.POST("", handler.CreateUser)
	group.GET("", handler.GetAll)
	group.GET("/:id", handler.GetByID)

	return handler
}

func ErrResponse(err error) map[string]string {
	if err == nil {
		return map[string]string{
			"error": "unknown error",
		}
	}

	return map[string]string{
		"error": err.Error(),
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	if user.Name == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	createdUser, err := h.userUseCase.CreateUser(user)
	if err == nil {
		return c.JSON(http.StatusCreated, createdUser)
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

func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	user, err := h.userUseCase.GetByID(int64(id))
	if err == nil {
		return c.JSON(http.StatusOK, user)
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(err))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}

func (h *UserHandler) GetAll(c echo.Context) error {

	users, err := h.userUseCase.GetAll()
	if err == nil {
		return c.JSON(http.StatusOK, users)
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrResponse(err))
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}
