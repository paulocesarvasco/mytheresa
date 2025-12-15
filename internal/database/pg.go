package database

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(user, password, dbname, port string) (db *gorm.DB, close func() error) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	log := logs.Logger()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log, logger.Config{
			LogLevel: logger.LogLevel(logs.CurrentLevel()),
		}),
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
