package logs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type customLogger struct {
	level  slog.LevelVar
	logger *slog.Logger
}

var defaultLogger customLogger

type ctxKey struct{}

var loggerKey ctxKey

func (l *customLogger) Info(ctx context.Context, msg string, args ...any) {
	from(ctx).InfoContext(ctx, msg, args...)
}

func (l *customLogger) Warn(ctx context.Context, msg string, args ...any) {
	from(ctx).WarnContext(ctx, msg, args...)
}

func (l *customLogger) Error(ctx context.Context, msg string, args ...any) {
	from(ctx).ErrorContext(ctx, msg, args...)
}

func (l *customLogger) Debug(ctx context.Context, msg string, args ...any) {
	from(ctx).DebugContext(ctx, msg, args...)
}

func (l *customLogger) Printf(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}

func Init(level slog.Level) ApiLogger {
	defaultLogger.level.Set(level)

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: &defaultLogger.level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	})
	defaultLogger.logger = slog.New(h)
	slog.SetDefault(defaultLogger.logger)
	return &defaultLogger
}

func CurrentLevel() slog.Level {
	return defaultLogger.level.Level()
}

func Logger() ApiLogger {
	if defaultLogger.logger == nil {
		Init(slog.LevelInfo)
	}
	return &defaultLogger
}

// Into attaches a logger to the context.
func into(ctx context.Context, l *slog.Logger) context.Context {
	if l == nil {
		l = defaultLogger.logger
	}
	return context.WithValue(ctx, loggerKey, l)
}

// From retrieves logger from context
func from(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return defaultLogger.logger
	}
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok && l != nil {
		return l
	}
	return defaultLogger.logger
}
