package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/athifirshad/eucalyptus/db"
	"github.com/athifirshad/eucalyptus/internal/data"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	config
	logger  *zap.Logger
	router  *chi.Mux
	models  data.Models //handmade queries
	queries *db.Queries //sqlc generated queries
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

	dbPool, err := openDB(cfg)
	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	defer dbPool.Close()
	defer logger.Sync()
	sugar := logger.Sugar()

	router := chi.NewRouter()

	app := &application{
		config: cfg,
		logger: logger,
		router: router,
		models: data.NewModels(dbPool),
		queries: db.New(dbPool), 
	}
	sugar.Infof("Database connection estabilished")
	sugar.Infof("Starting %s server on %s", cfg.env, cfg.port)
	app.router.Use(app.logRequest)
	app.Routes()
	http.ListenAndServe(cfg.port, app.router)

}

func openDB(cfg config) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}
