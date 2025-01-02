package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/codeandlearn1991/newsapi/internal/logger"
	"github.com/codeandlearn1991/newsapi/internal/news"
	"github.com/codeandlearn1991/newsapi/internal/postgres"
	"github.com/codeandlearn1991/newsapi/internal/router"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := postgres.NewDB(&postgres.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		DBName:   os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		User:     os.Getenv("DATABASE_USER"),
		Port:     os.Getenv("DATABASE_PORT"),
		SSLMode:  "disable",
	})
	if err != nil {
		log.Error("db error", "err", err)
		os.Exit(1)
	}
	newsStore := news.NewStore(db)

	r := router.New(newsStore)
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
