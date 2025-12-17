package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	catalogapi "github.com/mytheresa/go-hiring-challenge/internal/api/catalog"
	categoriesapi "github.com/mytheresa/go-hiring-challenge/internal/api/categories"
	"github.com/mytheresa/go-hiring-challenge/internal/api/middlewares"
	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	"github.com/mytheresa/go-hiring-challenge/internal/categories"
	"github.com/mytheresa/go-hiring-challenge/internal/database"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
)

func main() {
	logs.Init(slog.LevelDebug)

	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	// Signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	db, close := database.New(dsn)
	defer close()

	// Initialize dependencies
	productStore := repository.NewProductStore(db)
	catalogService := catalog.New(productStore)
	catalogHandler := catalogapi.New(catalogService)

	categoriesStore := repository.NewCategoryStore(db)
	categoriesService := categories.New(categoriesStore)
	categoriesHandler := categoriesapi.New(categoriesService)

	// Router
	r := chi.NewRouter()

	// Defaults middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	// Custom logger
	r.Use(middlewares.RequestLogger)

	// Routes
	r.Mount("/catalog", catalogapi.Routes(catalogHandler))
	r.Mount("/categories", categoriesapi.Routes(categoriesHandler))

	// HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
		Handler: r,
	}

	// Start server
	go func() {
		log.Printf("starting server on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %s", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped gracefully")
}
