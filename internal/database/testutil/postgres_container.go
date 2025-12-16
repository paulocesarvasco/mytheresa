package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func StartPostgresContainer(t *testing.T, ctx context.Context) (*gorm.DB, func()) {
	t.Helper()

	migrations, err := getMigrationScripts()
	if err != nil {
		t.Fatalf("failed to retrieve migration files: %v", err)
	}

	pg, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithOrderedInitScripts(migrations...),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForLog("database system is ready to accept connections"),
			).WithDeadline(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	terminate := func() {
		err = pg.Terminate(ctx)
		if err != nil {
			t.Fatalf("failed to terminate database container: %v", err)
		}

	}

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		terminate()
		t.Fatalf("failed to retrieve postgres dsn container: %v", err)
	}

	db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{
		Logger: logger.New(logs.Logger(), logger.Config{
			LogLevel: logger.LogLevel(logs.CurrentLevel()),
		}),
		TranslateError: true,
	})
	if err != nil {
		terminate()
		t.Fatalf("failed initialize gorm session container: %v", err)
	}

	return db, terminate
}

func getMigrationScripts() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, fmt.Errorf("could not find project root")
		}
		dir = parent
	}
	var migrationScripts []string

	migrationsDir := filepath.Join(dir, "migrations")
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			migrationScripts = append(migrationScripts, filepath.Join(migrationsDir, f.Name()))
		}
	}

	slices.Sort(migrationScripts)

	return migrationScripts, nil
}
