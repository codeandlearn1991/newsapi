package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/codeandlearn1991/newsapi/internal/logger"
	"github.com/codeandlearn1991/newsapi/internal/router"
	"github.com/codeandlearn1991/newsapi/internal/store"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	r := router.New(store.New())
	wrappedRouter := logger.AddLoggerMid(log, logger.LoggerMid(r))

	log.Info("server starting on port 8080")

	if err := http.ListenAndServe(":8080", wrappedRouter); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
