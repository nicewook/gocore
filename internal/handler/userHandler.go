package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/internal/middlewares"
	"github.com/nicewook/gocore/pkg/contextutil"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(e *echo.Echo, userUseCase domain.UserUseCase) *UserHandler {
	handler := &UserHandler{userUseCase: userUseCase}

	group := e.Group("/users")
	group.GET("", handler.GetAll, middlewares.AllowRoles(domain.RoleAdmin))
	group.GET("/:id", handler.GetByID, middlewares.AllowRoles(domain.RoleAdmin, domain.RoleUser))

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

func (h *UserHandler) GetByID(c echo.Context) error {
	req := new(domain.GetByIDRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	ctx := c.Request().Context()
	user, err := h.userUseCase.GetByID(ctx, req.ID)
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
	logger := contextutil.GetLogger(c.Request().Context())
	logger.Info("UserHandler:GetAll")

	// Authorization 헤더 확인
	authHeader := c.Request().Header.Get("Authorization")
	logger.Debug("Authorization header", "value", authHeader)

	ctx := c.Request().Context()
	// 인증 없이도 사용자 목록 반환
	users, err := h.userUseCase.GetAll(ctx)
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
