package database

import (
	"context"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() (db *gorm.DB, close func() error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	log := logs.Logger()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log, logger.Config{
			LogLevel: logger.LogLevel(logs.CurrentLevel()),
		}),
		TranslateError: true,
	})
	if err != nil {
		log.Error(context.Background(), "failed to connect database", "error", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error(context.Background(), "failed to get database connection", "error", err)
	}

	return db, sqlDB.Close
}
