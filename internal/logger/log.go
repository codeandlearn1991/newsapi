package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

type CtxKey struct{}

func CtxWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		return ctx
	}

	if ctxLog, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok && ctxLog == logger {
		return ctx
	}

	return context.WithValue(ctx, CtxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}

func AddLoggerMid(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loggerCtx := CtxWithLogger(r.Context(), logger)
		r = r.Clone(loggerCtx)
		next.ServeHTTP(w, r)
	}
}

func LoggerMid(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := FromContext(r.Context())
		l.Info("request", "path", r.URL.String())
		next.ServeHTTP(w, r)
	}
}
