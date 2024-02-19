package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/athifirshad/eucalyptus/db"
	"github.com/athifirshad/eucalyptus/internal/data"
	"github.com/athifirshad/eucalyptus/internal/mailer"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
type application struct {
	config
	logger *zap.Logger
	router *chi.Mux
	models data.Models //handmade queries
	sqlc   *db.Queries //sqlc generated queries
	mailer *mailer.Mailer
	wg     sync.WaitGroup
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

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "47a0bd37235fa1", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "9a0ad4d8cdadb7", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Eucalyptus <no-reply@eucalyptus.net>", "SMTP sender")
	flag.Parse()

	//TODO Sentry reporting
	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, _ := config.Build()
	// logger := zap.Must(zap.NewProduction())
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

	//router := chi.NewRouter()

	mailer, err := mailer.New("sandbox.smtp.mailtrap.io", 587, "47a0bd37235fa1", "9a0ad4d8cdadb7", "Eucalyptus <no-reply@eucalyptus.net>")
	if err != nil {
		logger.Fatal("Failed to create mailer", zap.Error(err))
	}

	app := &application{
		config: cfg,
		logger: logger,
		router: chi.NewRouter(),
		models: data.NewModels(dbPool),
		sqlc:   db.New(dbPool),
		mailer: mailer,
	}
	sugar.Info("Database connection estabilished")
	app.router.Use(app.logRequest)
	app.Routes()
	if err := app.serve(); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
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
