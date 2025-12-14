package logs

import (
	"context"
	"log/slog"
	"os"
)

type CustomLogger struct{}

func Init(level slog.Level) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	})
	slog.SetDefault(slog.New(h))
}

// Into attaches a logger to the context (used by middleware).
func Into(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, CustomLogger{}, l)
}

func NewLogger() *CustomLogger {
	return &CustomLogger{}
}

func (cl *CustomLogger) Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func (cl *CustomLogger) Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func (cl *CustomLogger) Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

func (cl *CustomLogger) Debug(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}
