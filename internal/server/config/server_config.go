package config

import (
	"os"
	"path/filepath"

	"github.com/driif/golang-test-task/internal/server/config/env"
	"github.com/driif/golang-test-task/pkg/tests"
)

// EchoServer represents a subset of echo's config relevant to the app server.
type GinServer struct {
	Debug                          bool
	ListenAddress                  string
	HideInternalServerErrorDetails bool
	BaseURL                        string
	EnableCORSMiddleware           bool
	EnableLoggerMiddleware         bool
	EnableRecoverMiddleware        bool
	EnableRequestIDMiddleware      bool
	EnableTrailingSlashMiddleware  bool
	EnableSecureMiddleware         bool
	EnableCacheControlMiddleware   bool
	SecureMiddleware               EchoServerSecureMiddleware
}

// EchoServerSecureMiddleware represents a subset of echo's secure middleware config relevant to the app server.
// https://github.com/labstack/echo/blob/master/middleware/secure.go
type EchoServerSecureMiddleware struct {
	XSSProtection         string
	ContentTypeNosniff    string
	XFrameOptions         string
	HSTSMaxAge            int
	HSTSExcludeSubdomains bool
	ContentSecurityPolicy string
	CSPReportOnly         bool
	HSTSPreloadEnabled    bool
	ReferrerPolicy        string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Rabbitmq struct {
	Addr     string
	Username string
	Password string
}

// Server represents the config of the Server relevant to the app server, containing all the other config structs.
type Server struct {
	Environment string
	Gin         GinServer
	Redis       Redis
	Rabbitmq    Rabbitmq
}

func DefaultServiceConfigFromEnv() Server {
	if !tests.RunningInTest() {
		env.DotEnvTryLoad(filepath.Join(env.GetProjectRootDir(), ".env"), os.Setenv)
	}

	return Server{
		Environment: env.GetEnv("ENVIRONMENT", "development"),

		Gin: GinServer{
			Debug:                          env.GetEnvAsBool("SERVER_ECHO_DEBUG", false),
			ListenAddress:                  env.GetEnv("SERVER_ECHO_LISTEN_ADDRESS", ":3000"),
			HideInternalServerErrorDetails: env.GetEnvAsBool("SERVER_ECHO_HIDE_INTERNAL_SERVER_ERROR_DETAILS", true),
			BaseURL:                        env.GetEnv("SERVER_ECHO_BASE_URL", "http://localhost:8080"),
			EnableCORSMiddleware:           env.GetEnvAsBool("SERVER_ECHO_ENABLE_CORS_MIDDLEWARE", true),
			EnableLoggerMiddleware:         env.GetEnvAsBool("SERVER_ECHO_ENABLE_LOGGER_MIDDLEWARE", true),
			EnableRecoverMiddleware:        env.GetEnvAsBool("SERVER_ECHO_ENABLE_RECOVER_MIDDLEWARE", true),
			EnableRequestIDMiddleware:      env.GetEnvAsBool("SERVER_ECHO_ENABLE_REQUEST_ID_MIDDLEWARE", true),
			EnableTrailingSlashMiddleware:  env.GetEnvAsBool("SERVER_ECHO_ENABLE_TRAILING_SLASH_MIDDLEWARE", true),
			EnableSecureMiddleware:         env.GetEnvAsBool("SERVER_ECHO_ENABLE_SECURE_MIDDLEWARE", true),
			EnableCacheControlMiddleware:   env.GetEnvAsBool("SERVER_ECHO_ENABLE_CACHE_CONTROL_MIDDLEWARE", true),
			// see https://echo.labstack.com/middleware/secure
			// see https://github.com/labstack/echo/blob/master/middleware/secure.go
			SecureMiddleware: EchoServerSecureMiddleware{
				XSSProtection:         env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_XSS_PROTECTION", "1; mode=block"),
				ContentTypeNosniff:    env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_CONTENT_TYPE_NOSNIFF", "nosniff"),
				XFrameOptions:         env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_X_FRAME_OPTIONS", "SAMEORIGIN"),
				HSTSMaxAge:            env.GetEnvAsInt("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_MAX_AGE", 0),
				HSTSExcludeSubdomains: env.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_EXCLUDE_SUBDOMAINS", false),
				ContentSecurityPolicy: env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_CONTENT_SECURITY_POLICY", ""),
				CSPReportOnly:         env.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_CSP_REPORT_ONLY", false),
				HSTSPreloadEnabled:    env.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_PRELOAD_ENABLED", false),
				ReferrerPolicy:        env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_REFERRER_POLICY", ""),
			},
		},

		Redis: Redis{
			Addr:     env.GetEnv("REDIS_ADDR", "localhost:6379"),
			Password: env.GetEnv("REDIS_PASSWORD", ""),
			DB:       env.GetEnvAsInt("REDIS_DB", 0),
		},

		Rabbitmq: Rabbitmq{
			Addr:     env.GetEnv("RABBITMQ_ADDR", "amqp://guest:guest@localhost:5672/"),
			Username: env.GetEnv("RABBITMQ_USERNAME", "guest"),
			Password: env.GetEnv("RABBITMQ_PASSWORD", "guest"),
		},
	}
}
