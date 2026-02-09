package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/validacion-pases/internal/app"
	"github.com/example/validacion-pases/internal/config"
	"github.com/example/validacion-pases/internal/platform/db"
	"github.com/example/validacion-pases/internal/platform/observability"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	ctx := context.Background()
	shutdownTelemetry, err := observability.InitTelemetry(ctx, cfg)
	if err != nil {
		logger.Error("telemetry initialization failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdownTelemetry(context.Background()); err != nil {
			logger.Error("telemetry shutdown failed", "error", err)
		}
	}()

	dbConn, err := db.NewMySQL(cfg)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer closeDB(logger, dbConn)

	handler, err := app.New(ctx, cfg, dbConn, logger)
	if err != nil {
		logger.Error("app initialization failed", "error", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	go func() {
		logger.Info("api server started", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server crashed", "error", err)
			os.Exit(1)
		}
	}()

	waitForShutdown(logger, srv, cfg)
}

func waitForShutdown(logger *slog.Logger, srv *http.Server, cfg config.Config) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	logger.Info("graceful shutdown started")

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}

	logger.Info("shutdown completed", "at", time.Now().UTC())
}

func closeDB(logger *slog.Logger, dbConn *sql.DB) {
	if err := dbConn.Close(); err != nil {
		logger.Error("database close failed", "error", err)
	}
}
