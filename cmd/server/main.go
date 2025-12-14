package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	catalogapi "github.com/mytheresa/go-hiring-challenge/internal/api/catalog"
	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	"github.com/mytheresa/go-hiring-challenge/internal/database"
	"github.com/mytheresa/go-hiring-challenge/internal/repository"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	// Signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize dependencies
	prodRepo := repository.New(db)
	catalogService := catalog.New(prodRepo)
	catalogHandler := catalogapi.New(catalogService)

	// Router
	r := chi.NewRouter()

	// Defaults middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/catalog", catalogHandler.GetProducts)

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
