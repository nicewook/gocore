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

	// Parse query parameters and bind to request struct
	req := new(domain.GetAllUsersRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	// Validate request parameters
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse(domain.ErrInvalidInput))
	}

	// Set default values if not provided
	if req.Limit == 0 {
		req.Limit = 10 // Default limit
	}

	// Authorization 헤더 확인 (로깅 목적)
	authHeader := c.Request().Header.Get("Authorization")
	logger.Debug("Authorization header", "value", authHeader)

	ctx := c.Request().Context()

	// 통합된 GetAll 메서드 호출 (페이지네이션 정보 포함)
	response, err := h.userUseCase.GetAll(ctx, req)
	if err == nil {
		return c.JSON(http.StatusOK, response)
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		// 사용자를 찾지 못한 경우 빈 배열과 함께 페이지네이션 정보 반환
		emptyResponse := &domain.GetAllResponse{
			Users:      []domain.User{},
			TotalCount: 0,
			Offset:     req.Offset,
			Limit:      req.Limit,
			HasMore:    false,
		}
		return c.JSON(http.StatusOK, emptyResponse)
	default:
		return c.JSON(http.StatusInternalServerError, ErrResponse(domain.ErrInternal))
	}
}
