package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"

	"github.com/nicewook/gocore/internal/config"
)

func RegisterMiddlewares(cfg *config.Config, e *echo.Echo) {

	// ✅ Trailing Slash 제거 및 301 리디렉트 설정
	e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently, // 301 리디렉트
	}))

	// ✅ RequestID: 각 요청에 고유한 ID 부여 (추적 및 디버깅 목적)
	e.Use(middleware.RequestID())

	// ✅ Logger: 요청 및 응답 로깅 설정
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${method} ${uri} ${status} request_id=${id}\n`,
	}))

	// ✅ Recover: 패닉 발생 시 복구 및 로그 출력
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 스택 크기: 1KB
		LogLevel:  log.ERROR,
	}))

	// ✅ Gzip: 응답 압축 (성능 최적화)
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:     5,   // 압축 레벨 (1-9)
		MinLength: 256, // 최소 압축 크기 (256바이트 이상만 압축)
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "metrics") // 특정 경로는 압축 제외
		},
	}))

	// ✅ BodyDump: 요청/응답 본문 로깅 (대형 요청 감지)
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		// 대형 요청/응답 감지 (1MB 초과 시 경고)
		if len(reqBody) > 1024*1024 {
			e.Logger.Warnf("Large request body: %d bytes for %s %s",
				len(reqBody), c.Request().Method, c.Path())
		}
		if len(resBody) > 1024*1024 {
			e.Logger.Warnf("Large response body: %d bytes for %s %s",
				len(resBody), c.Request().Method, c.Path())
		}

		// 민감 경로는 본문 로깅 제외
		sensitivePaths := []string{"/login", "/register"} // 예시
		for _, path := range sensitivePaths {
			if c.Path() == path {
				e.Logger.Infof("%s %s [Sensitive path, body logging skipped]",
					c.Request().Method, c.Path())
				return
			}
		}

		// 디버그 모드에서만 본문 로깅 (1KB로 출력 제한)
		if c.Echo().Logger.Level() == log.DEBUG {
			// 클로저: 본문 길이 제한 함수
			limitBody := func(body []byte, max int) string {
				if len(body) > max {
					return string(body[:max]) + " [TRUNCATED]"
				}
				return string(body)
			}

			// 클로저: Content-Type 검사 함수. 텍스트인 경우만 로깅하기 위함
			isTextContent := func(contentType string) bool {
				return strings.HasPrefix(contentType, "application/json") ||
					strings.HasPrefix(contentType, "text/") ||
					strings.HasPrefix(contentType, "application/xml") ||
					strings.HasPrefix(contentType, "application/x-www-form-urlencoded")
			}

			// Content-Type 가져오기 (요청 및 응답)
			reqContentType := c.Request().Header.Get(echo.HeaderContentType)
			resContentType := c.Response().Writer.Header().Get(echo.HeaderContentType)

			// Content-Type이 텍스트일 때만 로깅
			if isTextContent(reqContentType) && isTextContent(resContentType) {
				e.Logger.Debugf("Request: %s", limitBody(reqBody, 1000))
				e.Logger.Debugf("Response: %s", limitBody(resBody, 1000))
			} else {
				e.Logger.Debugf("Request and Response are non-text content. Skipping log.")
			}
		}
	}))

	// ✅ BodyLimit: 요청 크기 제한 (2MB)
	e.Use(middleware.BodyLimit("2M"))

	// ✅ 핸들러 실행 시간 제한
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second, // 핸들러 내부 실행 시간
	}))

	// ✅ 서버 자체 타임아웃 설정
	e.Server.ReadTimeout = 10 * time.Second  // 요청 읽기 타임아웃
	e.Server.WriteTimeout = 40 * time.Second // 응답 쓰기 타임아웃 (Handler Timeout보다 길게)
	e.Server.IdleTimeout = 120 * time.Second // 유휴 연결 타임아웃

	// 여기서부터는 보안과 관련한 미들웨어
	// ✅ Rate Limiter: 요청 속도 제한 (DDoS 방지)
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(10),  // 초당 10회 요청
				Burst:     30,              // 버스트 최대 30회
				ExpiresIn: 1 * time.Minute, // 1분 주기로 리셋
			},
		),

		// 클라이언트 식별: IP
		IdentifierExtractor: func(c echo.Context) (string, error) {
			ip := c.RealIP()
			if ip == "" {
				e.Logger.Warn("RateLimiter: Failed to extract client IP")
				return "", errors.New("unable to determine client IP")
			}
			return "ip:" + ip, nil
		},

		// IdentifierExtractor 실패 시 처리
		ErrorHandler: func(c echo.Context, err error) error {
			e.Logger.Errorf("RateLimiter: Identifier extraction failed: %v", err)
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Invalid client identifier",
			})
		},

		// 요청 초과 시 처리
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			e.Logger.Warnf("RateLimiter: Rate limit exceeded for identifier: %s", identifier)
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Rate limit exceeded",
			})
		},
	}))

	// ✅ CORS: Cross-Origin Resource Sharing 설정
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Secure.CORSAllowOrigins, // 각 환경의 설정에 허용하는 도메인을 설정해둔다
		AllowMethods: []string{ // 허용할 HTTP 메서드
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.PATCH,
			echo.OPTIONS,
		},
		AllowHeaders: []string{ // 허용할 요청 헤더
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXCSRFToken,
		},
		AllowCredentials: true,                            // 쿠키 및 인증정보 포함 허용 (JWT 쿠키 기반 인증 시 필수)
		ExposeHeaders:    []string{echo.HeaderXRequestID}, // 클라이언트에 노출할 응답 헤더
		MaxAge:           86400,                           // 사전 요청(Preflight) 결과 캐싱 시간 (초 단위, 24시간)
	}))

	// ✅ CSRF: Cross-Site Request Forgery 방어
	if cfg.App.Env != "dev" {
		// CSRF token route handler
		e.GET("/csrf-token", func(c echo.Context) error {
			token := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
			return c.JSON(http.StatusOK, map[string]string{
				"csrf_token": token,
			})
		})

		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "header:" + echo.HeaderXCSRFToken,
			CookieSecure:   false,                // HTTPS에서만 쿠키 전송
			CookiePath:     "/",                  // 이 설정 추가
			CookieName:     "_csrf",              // 이 설정 추가
			CookieHTTPOnly: true,                 // JavaScript 접근 금지
			CookieSameSite: http.SameSiteLaxMode, // 동일 출처 외 요청 차단
		}))
	}

	// ✅ Secure Headers: 다양한 보안 헤더 설정
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",

		// HSTS
		HSTSMaxAge:            15768000, // 6개월
		HSTSPreloadEnabled:    true,
		HSTSExcludeSubdomains: false, // 서브도메인까지 HTTPS 강제 적용

		// Content Security Policy (CSP)
		ContentSecurityPolicy: "default-src 'self'; frame-ancestors 'none'; form-action 'self';" +
			"permissions-policy: camera=(), microphone=(), geolocation=();",

		// CSP 테스트용 (필요 시)
		CSPReportOnly: false,

		// Referrer-Policy 강화
		ReferrerPolicy: "no-referrer",
	}))
}
