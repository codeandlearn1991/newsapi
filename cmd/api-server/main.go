package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/codeandlearn1991/newsapi/internal/logger"
	"github.com/codeandlearn1991/newsapi/internal/router"
	"github.com/codeandlearn1991/newsapi/internal/store"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	r := router.New(store.New())
	wrappedRouter := logger.AddLoggerMid(log, logger.Middleware(r))

	log.Info("server starting on port 8080")

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           wrappedRouter,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
