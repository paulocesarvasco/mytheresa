package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

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
	// Signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		logs.Logger().Error(ctx, "error loading .env file", "error", err)
		os.Exit(1)
	}

	log := logs.Init()

	// Initialize database
	db, close := database.New()
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{fmt.Sprintf("http://localhost:%s", os.Getenv("SWAGGER_PORT"))},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	// Custom logger
	r.Use(middlewares.RequestLogger)

	// Routes
	r.Mount("/catalog", catalogapi.Routes(catalogHandler))
	r.Mount("/categories", categoriesapi.Routes(categoriesHandler))

	// HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("HTTP_HOST"), os.Getenv("HTTP_PORT")),
		Handler: r,
	}

	// Start server
	go func() {
		log.Info(ctx, "starting server", "address", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(ctx, "server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Info(ctx, "shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error(ctx, "server shutdown error", "error", err)
		os.Exit(1)
	}

	log.Info(ctx, "server stopped gracefully")
}
