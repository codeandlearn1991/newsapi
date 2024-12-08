package logger_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/codeandlearn1991/newsapi/internal/logger"
)

func Test_CtxWithLogger(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    context.Context
		logger *slog.Logger
		exists bool
	}{
		{
			name: "returns context without logger",
			ctx:  context.Background(),
		},
		{
			name:   "return ctx as it is",
			ctx:    context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			exists: true,
		},
		{
			name:   "inject logger",
			ctx:    context.Background(),
			logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
			exists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := logger.CtxWithLogger(tc.ctx, tc.logger)

			_, ok := ctx.Value(logger.CtxKey{}).(*slog.Logger)
			if tc.exists != ok {
				t.Errorf("expected :%v got: %v", tc.exists, ok)
			}
		})
	}
}

func Test_FromContext(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		expected bool
	}{
		{
			name:     "logger exists",
			ctx:      context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			expected: true,
		},
		{
			name:     "new logger returned",
			ctx:      context.Background(),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logger.FromContext(tc.ctx)

			if tc.expected && logger == nil {
				t.Errorf("expected: %v, got: %v", tc.expected, logger)
			}
		})
	}
}
