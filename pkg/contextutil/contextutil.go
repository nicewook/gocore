package contextutil

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type contextKey string

var loggerContextKey contextKey = "logger_context_key"

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*slog.Logger); ok {
		return logger
	}

	logger := slog.Default()
	logger.Error("Default logger used",
		slog.String("reason", "no logger found in context"),
	)
	return logger
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, echo.HeaderXRequestID, requestID)
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(echo.HeaderXRequestID).(string); ok {
		return requestID
	}
	return "no-request-id-in-context"
}
