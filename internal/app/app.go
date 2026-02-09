package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/example/validacion-pases/internal/config"
	"github.com/example/validacion-pases/internal/repository/mysql"
	"github.com/example/validacion-pases/internal/security/auth"
	secheaders "github.com/example/validacion-pases/internal/security/headers"
	"github.com/example/validacion-pases/internal/transport/http/handlers"
	"github.com/example/validacion-pases/internal/transport/http/middleware"
	"github.com/example/validacion-pases/internal/usecase"
)

func New(ctx context.Context, cfg config.Config, db *sql.DB, logger *slog.Logger) (http.Handler, error) {
	validator, err := auth.NewJWTValidator(ctx, cfg.JWTAlg, cfg.JWTIssuer, cfg.JWTAudience, cfg.JWTClockSkew, cfg.JWTHSSecret, cfg.JWKSURL, cfg.JWTRefresh)
	if err != nil {
		return nil, err
	}

	var tokenSvc *usecase.TokenService
	issuer, issueErr := auth.NewTokenIssuer(cfg.JWTAlg, cfg.JWTIssuer, cfg.JWTAudience, cfg.JWTTokenTTL, cfg.JWTHSSecret)
	if issueErr == nil {
		userStore := auth.NewUserStore(cfg.TokenUsers)
		tokenSvc = usecase.NewTokenService(userStore, issuer)
	} else {
		logger.Warn("token issuance disabled", "reason", issueErr.Error())
		tokenSvc = usecase.NewTokenService(nil, nil)
	}

	repo := mysql.NewRecordRepository(db)
	svc := usecase.NewRecordService(repo)
	health := handlers.NewHealthHandler(db)
	records := handlers.NewRecordHandler(svc)
	tokenHandler := handlers.NewTokenHandler(tokenSvc)

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestLog(logger))
	r.Use(chimiddleware.Timeout(cfg.RequestTimeout))
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Heartbeat("/ping"))
	r.Use(chimiddleware.Compress(5))
	r.Use(chimiddleware.Throttle(100))
	r.Use(httprate.LimitByIP(cfg.RateLimitRequests, cfg.RateLimitWindow))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(secheaders.Middleware)
	r.Use(bodyLimit(cfg.BodyLimitBytes))

	reg := prometheus.NewRegistry()
	metrics := middleware.NewMetrics(reg)
	r.Use(metrics.Middleware)

	r.Get("/healthz", health.Liveness)
	r.Get("/readyz", health.Readiness)
	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	r.Route("/v1", func(v1 chi.Router) {
		v1.Use(chimiddleware.AllowContentType("application/json"))
		v1.Post("/token", tokenHandler.Issue)
		v1.With(middleware.AuthBearer(validator)).Post("/records", records.Create)
	})

	wrapped := otelhttp.NewHandler(r, "http.server", otelhttp.WithSpanNameFormatter(func(_ string, req *http.Request) string {
		return req.Method + " " + req.URL.Path
	}))

	return wrapped, nil
}

func bodyLimit(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
