package main

import (
	"context"
	"database/sql"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
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
	config
	logger *zap.Logger
	router *chi.Mux
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
	var cfg config
	flag.StringVar(&cfg.port, "port", "localhost:4000", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("NEON_DSN"), "PostgreSQL DSN")
	flag.Parse()

	logger := zap.Must(zap.NewProduction())
	if cfg.env == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	defer db.Close()
	defer logger.Sync()
	sugar := logger.Sugar()

	router := chi.NewRouter()

	app := &application{
		config: cfg,
		logger: logger,
		router: router,
	}
	sugar.Infof("Neon database connection estabilished")
	sugar.Infof("Starting %s server on %s", cfg.env, cfg.port)
	app.router.Use(app.logRequest)
	app.Routes()
	http.ListenAndServe(cfg.port, app.router)

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}