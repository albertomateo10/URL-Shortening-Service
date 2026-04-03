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

	"github.com/albertomateo10/url-shortener/backend/internal/config"
	"github.com/albertomateo10/url-shortener/backend/internal/handler"
	"github.com/albertomateo10/url-shortener/backend/internal/middleware"
	"github.com/albertomateo10/url-shortener/backend/internal/repository"
	"github.com/albertomateo10/url-shortener/backend/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database connection pool
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to PostgreSQL")

	// Run migrations
	if err := runMigrations(dbPool, cfg.DatabaseURL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("connected to Redis")

	// Repositories
	urlRepo := repository.NewURLRepository(dbPool)
	clickRepo := repository.NewClickRepository(dbPool)

	// Services
	cacheSvc := service.NewCacheService(rdb)
	urlSvc := service.NewURLService(urlRepo, cacheSvc, cfg.BaseURL)
	analyticsSvc := service.NewAnalyticsService(clickRepo, urlRepo)
	clickLogger := service.NewClickLogger(clickRepo, urlRepo)

	// Geo service (optional — only active when GEODB_PATH is set)
	geoSvc := service.NewGeoService(cfg.GeoDBPath)
	defer geoSvc.Close()

	// Handlers
	urlHandler := handler.NewURLHandler(urlSvc)
	redirectHandler := handler.NewRedirectHandler(urlSvc, clickLogger, geoSvc)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS(cfg.FrontendURL))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/urls", urlHandler.Create)
		r.Get("/urls", urlHandler.List)
		r.Get("/urls/{id}", urlHandler.Get)
		r.Delete("/urls/{id}", urlHandler.Delete)

		r.Get("/urls/{id}/analytics/clicks", analyticsHandler.GetClicksOverTime)
		r.Get("/urls/{id}/analytics/sources", analyticsHandler.GetSources)
	})

	// Redirect route
	r.Get("/r/{shortCode}", redirectHandler.Redirect)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	log.Printf("server starting on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations(pool *pgxpool.Pool, dbURL string) error {
	ctx := context.Background()

	// Read and execute migration files directly
	migrations := []string{
		"migrations/001_create_urls.up.sql",
		"migrations/002_create_click_events.up.sql",
	}

	for _, file := range migrations {
		sql, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			// Ignore errors for already-existing tables
			log.Printf("migration %s: %v (may already be applied)", file, err)
		}
	}

	log.Println("migrations applied")
	return nil
}
