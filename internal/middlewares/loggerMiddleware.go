package middlewares

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/nicewook/gocore/pkg/contextutil"
)

func LoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{

		BeforeNextFunc: func(c echo.Context) {
			// 요청이 시작될 때 실행되는 함수.
			// - 요청 컨텍스트에 `request_id`가 포함된 로거를 추가하여 이후 handler, usecase, repository 등에서 로깅할 때 사용.
			req := c.Request()
			requestID := contextutil.GetRequestID(c.Request().Context())
			ctxLogger := logger.With(slog.String("request_id", requestID))
			ctx := contextutil.WithLogger(req.Context(), ctxLogger)
			c.SetRequest(req.WithContext(ctx))
		},

		// 요청 및 응답에서 로깅할 값들
		LogRequestID: true, // 요청 ID 로깅
		LogStatus:    true, // HTTP 응답 상태 코드 로깅
		LogMethod:    true, // HTTP 메서드 (GET, POST 등) 로깅
		LogURIPath:   true, // 요청된 URI 경로 로깅
		LogRemoteIP:  true, // 클라이언트 IP 주소 로깅
		LogUserAgent: true, // 요청한 클라이언트의 User-Agent 로깅
		LogReferer:   true, // Referer 헤더(어디서 요청이 왔는지) 로깅
		LogLatency:   true, // 요청 완료까지 걸린 시간 로깅
		LogError:     true, // 에러 발생 시 로깅
		HandleError:  true, // 에러 발생 시 글로벌 에러 핸들러로 전달하여 적절한 응답을 반환

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// 공통 로깅 필드 설정 (모든 요청에서 공통적으로 기록할 항목)
			baseLogger := logger.With(
				slog.String("request_id", v.RequestID),
				slog.Int("status", v.Status),
				slog.String("method", v.Method),
				slog.String("path", v.URIPath),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("user_agent", v.UserAgent),
				slog.String("referer", v.Referer),
				slog.String("latency", v.Latency.String()),
			)
			if v.Error != nil {
				baseLogger.With(slog.String("err", v.Error.Error())).Error("REQUEST_ERROR")
			} else {
				baseLogger.Info("REQUEST_SUCCESS")
			}
			return nil
		},
	})
}
