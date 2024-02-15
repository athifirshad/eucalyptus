package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/athifirshad/eucalyptus/internal/data"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}
type application struct {
	config  config
	logger  *zap.Logger
	router  *chi.Mux
	queries *data.Queries
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Received request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)

		next.ServeHTTP(w, r)
	})
}

func main() {
	banner := `
	   ______  __________   ____  _____  ________  ______
	  / __/ / / / ___/ _ | / /\ \/ / _ \/_  __/ / / / __/
	 / _// /_/ / /__/ __ |/ /__\  / ___/ / / / /_/ /\ \  
	/___/\____/\___/_/ |_/____//_/_/    /_/  \____/___/  ` + "\n"
	fmt.Print(banner)
	var cfg config
	flag.StringVar(&cfg.port, "port", "localhost:4000", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("NEON_DSN"), "PostgreSQL DSN")
	flag.Parse()

	//TODO Sentry reporting

	logger := zap.Must(zap.NewProduction())
	if cfg.env == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	pool, err := openDB(cfg)
	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	defer logger.Sync()
	sugar := logger.Sugar()

	router := chi.NewRouter()

	app := &application{
		config:  cfg,
		logger:  logger,
		router:  router,
		queries: data.New(pool),
	}

	sugar.Infof("Database connection estabilished")
	sugar.Infof("Starting %s server on %s", cfg.env, cfg.port)
	app.queries = data.New(pool)
	app.router.Use(app.logRequest)
	app.Routes()
	http.ListenAndServe(cfg.port, app.router)

}

func openDB(cfg config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database configuration: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), poolConfig.ConnString())
	if err != nil {
		fmt.Errorf("unable to connect to database: %w", err)
		os.Exit(1)
	}
	defer pool.Close()

	return pool, nil
}
