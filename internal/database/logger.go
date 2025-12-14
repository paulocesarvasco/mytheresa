// Package database provides database-related infrastructure components.
//
// This package includes an exclusive implementation of a custom GORM logger,
// designed specifically to satisfy gorm.io/gorm/logger.Interface and to
// enhance observability of database operations such as query execution,
// slow queries, and database errors.
//
// The custom logger integrates with the application’s logging configuration
// while remaining strictly scoped to GORM usage, ensuring that ORM-specific
// logging concerns do not leak into API or service-layer code.
package database

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"gorm.io/gorm/logger"
)

type GormCustomLogger struct {
	logs.CustomLogger
	level         logger.LogLevel
	slowThreshold time.Duration
}

func NewLogger() *GormCustomLogger {
	logs.Init(slog.LevelDebug)
	return &GormCustomLogger{}
}

func (gl *GormCustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	cp := *gl
	cp.level = level
	return &cp
}

func (gl *GormCustomLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if gl.level == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	attrs := []any{
		"sql", sql,
		"rows", rows,
		"duration_ms", elapsed.Milliseconds(),
	}

	if err != nil && !errors.Is(err, logger.ErrRecordNotFound) {
		if gl.level >= logger.Error {
			gl.Error(ctx, "query failed", append(attrs, "err", err)...)
		}
		return
	}

	if gl.slowThreshold > 0 && elapsed >= gl.slowThreshold {
		if gl.level >= logger.Warn {
			gl.Warn(ctx, "slow query", attrs...)
		}
		return
	}

	if gl.level >= logger.Info {
		gl.Info(ctx, "query executed", attrs...)
	}
}
