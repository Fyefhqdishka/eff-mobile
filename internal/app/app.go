package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/Fyefhqdishka/eff-mobile/docs"
	"github.com/Fyefhqdishka/eff-mobile/internal/client"
	"github.com/Fyefhqdishka/eff-mobile/internal/config"
	"github.com/Fyefhqdishka/eff-mobile/internal/handlers"
	service "github.com/Fyefhqdishka/eff-mobile/internal/service"
	"github.com/Fyefhqdishka/eff-mobile/internal/storage"
	"github.com/Fyefhqdishka/eff-mobile/internal/storage/repositories"
	"github.com/Fyefhqdishka/eff-mobile/pkg/routes"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type App struct {
	db     *sql.DB
	server *http.Server
}

func (s *App) Run() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %v", err)
	}
	return nil
}

func (s *App) Stop() error {
	var errs []error

	if err := s.db.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close database: %v", err))
	}

	if err := s.server.Shutdown(context.Background()); err != nil {
		errs = append(errs, fmt.Errorf("failed to shutdown server: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// New creates new instance of application, sets the dependencies and applies migrations
func New(cfg *config.Config) (*App, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := storage.ConnectDB(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log := initLogging()

	baseURL := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	client := client.NewClient(baseURL, log)

	storage := repositories.NewSongRepository(db, log)

	service := service.NewService(storage, client, log)

	h := *handlers.NewHandlers(log, service)

	r := mux.NewRouter()
	routes.RegisterRoutes(r, h)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	log.Info("server starting in port", cfg.Server.Port)

	app := &App{
		db: db,
		server: &http.Server{
			Addr:         addr,
			Handler:      r,
			WriteTimeout: cfg.Server.Timeout,
			ReadTimeout:  cfg.Server.Timeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}

	return app, nil
}

func initLogging() *slog.Logger {
	logFileName := "logs/app-" + time.Now().Format("2006-01-02") + ".log"
	logfile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Не удалось открыть файл для логов", "error", err)
		os.Exit(1)
	}

	handler := slog.NewTextHandler(logfile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}
