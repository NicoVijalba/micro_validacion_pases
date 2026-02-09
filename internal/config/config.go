// Package config loads and validates runtime configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains application settings loaded from the environment.
type Config struct {
	ServiceName       string
	Environment       string
	HTTPAddr          string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	RequestTimeout    time.Duration
	ShutdownTimeout   time.Duration
	BodyLimitBytes    int64
	LogLevel          slog.Level

	DBDSN             string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	AuthMode     string
	JWTAlg       string
	JWTIssuer    string
	JWTAudience  string
	JWTClockSkew time.Duration
	JWKSURL      string
	JWTRefresh   time.Duration
	JWTHSSecret  string
	JWTTokenTTL  time.Duration
	TokenUsers   map[string]string

	RateLimitRequests int
	RateLimitWindow   time.Duration
	AllowedOrigins    []string

	OTelEnabled  bool
	OTelEndpoint string
	OTelInsecure bool
}

// Load builds the application configuration from environment variables and validates core constraints.
func Load() (Config, error) {
	cfg := Config{
		ServiceName:       getEnv("SERVICE_NAME", "validacion-pases"),
		Environment:       getEnv("ENV", "dev"),
		HTTPAddr:          getEnv("HTTP_ADDR", ":8080"),
		ReadTimeout:       mustDuration("HTTP_READ_TIMEOUT", "10s"),
		ReadHeaderTimeout: mustDuration("HTTP_READ_HEADER_TIMEOUT", "5s"),
		WriteTimeout:      mustDuration("HTTP_WRITE_TIMEOUT", "15s"),
		IdleTimeout:       mustDuration("HTTP_IDLE_TIMEOUT", "60s"),
		RequestTimeout:    mustDuration("HTTP_REQUEST_TIMEOUT", "8s"),
		ShutdownTimeout:   mustDuration("HTTP_SHUTDOWN_TIMEOUT", "15s"),
		BodyLimitBytes:    mustInt64("HTTP_BODY_LIMIT_BYTES", 1048576),
		LogLevel:          parseLogLevel(getEnv("LOG_LEVEL", "INFO")),

		DBDSN:             getEnv("DB_DSN", "app:app@tcp(localhost:3306)/validacion_pases?parseTime=true&tls=false"),
		DBMaxOpenConns:    mustInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    mustInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: mustDuration("DB_CONN_MAX_LIFETIME", "30m"),
		DBConnMaxIdleTime: mustDuration("DB_CONN_MAX_IDLE_TIME", "5m"),

		AuthMode:     getEnv("AUTH_MODE", "jwt"),
		JWTAlg:       strings.ToUpper(getEnv("JWT_ALG", "HS256")),
		JWTIssuer:    getEnv("JWT_ISSUER", "https://issuer.example.com"),
		JWTAudience:  getEnv("JWT_AUDIENCE", "validacion-pases"),
		JWTClockSkew: mustDuration("JWT_CLOCK_SKEW", "30s"),
		JWKSURL:      getEnv("JWT_JWKS_URL", ""),
		JWTRefresh:   mustDuration("JWT_REFRESH_INTERVAL", "5m"),
		JWTHSSecret:  getEnv("JWT_HS_SECRET", ""),
		JWTTokenTTL:  mustDuration("JWT_TOKEN_TTL", "1h"),
		TokenUsers:   parseTokenUsers(getEnv("TOKEN_USERS", "apiuser:change-me")),

		RateLimitRequests: mustInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   mustDuration("RATE_LIMIT_WINDOW", "1m"),
		AllowedOrigins:    splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),

		OTelEnabled:  mustBool("OTEL_ENABLED", false),
		OTelEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4318"),
		OTelInsecure: mustBool("OTEL_EXPORTER_OTLP_INSECURE", true),
	}

	if cfg.AuthMode != "jwt" {
		return Config{}, errors.New("only AUTH_MODE=jwt is implemented")
	}
	if cfg.JWTAlg == "RS256" && cfg.JWKSURL == "" {
		return Config{}, errors.New("JWT_JWKS_URL is required when JWT_ALG=RS256")
	}
	if cfg.JWTAlg == "HS256" && cfg.JWTHSSecret == "" {
		return Config{}, errors.New("JWT_HS_SECRET is required when JWT_ALG=HS256")
	}
	if len(cfg.TokenUsers) == 0 {
		return Config{}, errors.New("TOKEN_USERS must include at least one user:password pair")
	}

	return cfg, nil
}

func parseTokenUsers(raw string) map[string]string {
	users := make(map[string]string)
	entries := splitCSV(raw)
	for _, entry := range entries {
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			continue
		}
		user := strings.TrimSpace(parts[0])
		pass := strings.TrimSpace(parts[1])
		if user == "" || pass == "" {
			continue
		}
		users[user] = pass
	}
	return users
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustDuration(key, fallback string) time.Duration {
	raw := getEnv(key, fallback)
	d, err := time.ParseDuration(raw)
	if err != nil {
		panic(fmt.Sprintf("invalid duration for %s: %v", key, err))
	}
	return d
}

func mustInt(key string, fallback int) int {
	raw := getEnv(key, strconv.Itoa(fallback))
	v, err := strconv.Atoi(raw)
	if err != nil {
		panic(fmt.Sprintf("invalid int for %s: %v", key, err))
	}
	return v
}

func mustInt64(key string, fallback int64) int64 {
	raw := getEnv(key, strconv.FormatInt(fallback, 10))
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid int64 for %s: %v", key, err))
	}
	return v
}

func mustBool(key string, fallback bool) bool {
	raw := getEnv(key, strconv.FormatBool(fallback))
	v, err := strconv.ParseBool(raw)
	if err != nil {
		panic(fmt.Sprintf("invalid bool for %s: %v", key, err))
	}
	return v
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
