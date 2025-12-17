package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"

	"github.com/mytheresa/go-hiring-challenge/internal/database"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
)

func main() {
	log := logs.Init(slog.LevelInfo)

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Error(context.Background(), "Error loading .env file", "error", err)
		return
	}

	// Initialize database connection
	db, close := database.New()
	defer close()

	dir := os.Getenv("POSTGRES_SQL_DIR")
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Error(context.Background(), "reading directory failed", "error", err)
		return
	}

	// Filter and sort .sql files
	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, file := range sqlFiles {
		path := filepath.Join(dir, file.Name())

		content, err := os.ReadFile(path)
		if err != nil {
			log.Error(context.Background(), "reading file failed", "file", file.Name(), "error", err)
			return
		}

		sql := string(content)
		if err := db.Exec(sql).Error; err != nil {
			log.Error(context.Background(), "migration execution failed", "file", file.Name(), "error", err)
			return
		}

		log.Info(context.Background(), "executed successfully", "file", file.Name())
	}
}
