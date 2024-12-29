package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

// CtxKey for the logger.
type CtxKey struct{}

// CtxWithLogger returns the context enriched with the logger.
func CtxWithLogger(ctx context.Context, logger *slog.Logger) context.Context { //nolint:ireturn // needs to returns a enriched context.
	if logger == nil {
		return ctx
	}

	if ctxLog, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok && ctxLog == logger {
		return ctx
	}

	return context.WithValue(ctx, CtxKey{}, logger)
}

// FromContext fetch the logger from the context.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}

// AddLoggerMid adds the logger to the request context.
func AddLoggerMid(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loggerCtx := CtxWithLogger(r.Context(), logger)
		r = r.Clone(loggerCtx)
		next.ServeHTTP(w, r)
	}
}

// Middleware request logging middleware.
func Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := FromContext(r.Context())
		l.Info("request", "path", r.URL.String())
		next.ServeHTTP(w, r)
	}
}
