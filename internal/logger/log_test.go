package logger_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/codeandlearn1991/newsapi/internal/logger"
	"github.com/stretchr/testify/assert"
)

func Test_CtxWithLogger(t *testing.T) {
	testCases := []struct {
		name   string
		logger *slog.Logger
		key    logger.CtxKey
		value  *slog.Logger
		exists bool
	}{
		{
			name: "returns context without logger",
		},
		{
			name:   "return ctx as it is",
			key:    logger.CtxKey{},
			value:  slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
			exists: true,
		},
		{
			name:   "inject logger",
			logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
			exists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			if tc.value != nil {
				ctx = context.WithValue(ctx, tc.key, tc.value)
			}

			// Act
			ctx = logger.CtxWithLogger(ctx, tc.logger)

			// Assert
			_, ok := ctx.Value(logger.CtxKey{}).(*slog.Logger)
			assert.Equal(t, tc.exists, ok)
		})
	}
}

func Test_FromContext(t *testing.T) {
	testCases := []struct {
		name     string
		key      logger.CtxKey
		value    *slog.Logger
		expected bool
	}{
		{
			name:  "logger exists",
			key:   logger.CtxKey{},
			value: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
		},
		{
			name: "new logger returned",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			if tc.value != nil {
				ctx = context.WithValue(ctx, tc.key, tc.value)
			}

			// Act
			log := logger.FromContext(ctx)

			// Assert
			assert.True(t, log != nil)
		})
	}
}
