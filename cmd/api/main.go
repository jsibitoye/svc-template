package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jsibitoye/svc-template/internal/httpx"
	"github.com/jsibitoye/svc-template/internal/version"
)

func main() {
	port := getenv("PORT", "8080")
	shutdownSeconds := getenvInt("SHUTDOWN_TIMEOUT_SECONDS", 10)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("starting api",
		"version", version.Version,
		"commit", version.GitCommit,
		"built", version.BuildDate,
		"port", port,
	)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      httpx.NewMux(logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Info("shutting down")
	c, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownSeconds)*time.Second)
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
		_ = srv.Close()
	}
	logger.Info("stopped")
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
