package middlewares

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"

	"github.com/nicewook/gocore/internal/domain"
	"github.com/nicewook/gocore/pkg/contextutil"
)

// AllowRoles 는 허용된 역할만 접근할 수 있도록 하는 미들웨어이다.
func AllowRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 허용하는 역할중에 domain.RolePublic 역할이 있으면 토큰 검증을 하지 않는다.
			for _, role := range allowedRoles {
				if role == domain.RolePublic {
					return next(c)
				}
			}

			// token 속 사용자 roles 를 가져온다
			_, _, rolesInToken, err := contextutil.TokenToUser(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}

			// allowedRoles(허용하는 역할) 중에서 하나라도 일치하는 역할이 있으면 패스
			for _, role := range allowedRoles {
				if slices.Contains(rolesInToken, role) {
					return next(c)
				}
			}

			// 허용하는 역할중에 일치하는 역할이 없으면 403 에러 반환
			errMessage := fmt.Sprintf("insufficient permissions to access this resource. allowed: %v, user roles: %v", allowedRoles, rolesInToken)
			return echo.NewHTTPError(http.StatusForbidden, errMessage)
		}
	}
}
