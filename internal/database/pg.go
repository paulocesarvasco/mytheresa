package database

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(user, password, dbname, port string) (db *gorm.DB, close func() error) {
	log := NewLogger()

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log,
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
