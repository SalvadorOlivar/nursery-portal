package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	nurseryhttp "github.com/tuusuario/nursery-portal/internal/adapters/http"
	"github.com/tuusuario/nursery-portal/internal/adapters/repository/postgres"
	"github.com/tuusuario/nursery-portal/internal/application/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// dbURL := getEnv("DATABASE_URL", "postgres://nursery:nursery_dev@localhost:5432/nursery_portal?sslmode=disable")
	dbURL := getEnv("DATABASE_URL", "postgresql://postgres:vUbnXYya9Wdjcb1A@db.zeiucxhkmxngysemqyrn.supabase.co:5432/postgres")
	port := getEnv("PORT", "8080")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := runMigrations(ctx, dbURL); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("migrations applied successfully")

	authRepo := postgres.NewAuthRepository(pool)
	authSvc := services.NewAuthService(authRepo)
	if err := authSvc.EnsureAdmin(ctx, os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD")); err != nil {
		slog.Error("failed to ensure admin user", "error", err)
		os.Exit(1)
	}
	authHandler := nurseryhttp.NewAuthHandler(authSvc)
	authMiddleware := nurseryhttp.NewAuthMiddleware(authSvc)

	employeeRepo := postgres.NewEmployeeRepository(pool)
	employeeSvc := services.NewEmployeeService(employeeRepo, authSvc)
	employeeHandler := nurseryhttp.NewEmployeeHandler(employeeSvc)

	planifRepo := postgres.NewPlanificacionRepository(pool)
	turnoRepo := postgres.NewTurnoRepository(pool)
	dotacionRepo := postgres.NewDotacionRepository(pool)
	leaveRepo := postgres.NewLeaveRequestRepository(pool)
	compRepo := postgres.NewCompensatoryDayRepository(pool)
	planifSvc := services.NewPlanificacionService(planifRepo, turnoRepo, dotacionRepo, dotacionRepo, employeeRepo, leaveRepo, compRepo)
	planifHandler := nurseryhttp.NewPlanificacionHandler(planifSvc, employeeSvc)

	ausenciaSvc := services.NewAusenciaService(leaveRepo, compRepo)
	ausenciaHandler := nurseryhttp.NewAusenciaHandler(ausenciaSvc)

	intercambioRepo := postgres.NewIntercambioRepository(pool)
	intercambioSvc := services.NewIntercambioService(intercambioRepo, turnoRepo)
	intercambioHandler := nurseryhttp.NewIntercambioHandler(intercambioSvc)

	router := nurseryhttp.NewRouter(authHandler, authMiddleware, employeeHandler, planifHandler, ausenciaHandler, intercambioHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func runMigrations(ctx context.Context, dbURL string) error {
	db, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	return goose.UpContext(ctx, db, "migrations")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
