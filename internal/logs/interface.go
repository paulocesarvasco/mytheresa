package logs

import "context"

type ApiLogger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Printf(msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
}
