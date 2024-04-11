package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/athifirshad/eucalyptus/db"
	"github.com/athifirshad/eucalyptus/internal/data"
	"github.com/athifirshad/eucalyptus/internal/mailer"
	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
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
	redis struct {
		address string
	}
}
type application struct {
	config
	logger      *zap.Logger
	router      *chi.Mux
	models      data.Models //handmade queries
	sqlc        *db.Queries //sqlc generated queries
	mailer      *mailer.Mailer
	wg          sync.WaitGroup
	asynqClient *asynq.Client
}

func main() {
	banner := `
	   ______  __________   ____  _____  ________  ______
	  / __/ / / / ___/ _ | / /\ \/ / _ \/_  __/ / / / __/
	 / _// /_/ / /__/ __ |/ /__\  / ___/ / / / /_/ /\ \  
	/___/\____/\___/_/ |_/____//_/_/    /_/  \____/___/ ` + "\n"
	fmt.Print(banner)
	var cfg config

	// flag.StringVar(&cfg.port, "port", os.Getenv("PORT"), "API server port")

	cfg.port = ":" + os.Getenv("PORT")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Eucalyptus <no-reply@eucalyptus.net>", "SMTP sender")
	flag.StringVar(&cfg.redis.address, "redis-address", os.Getenv("REDISCLOUD_URL"), "Redis address")
	flag.Parse()

	//TODO Sentry reporting

	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	//logger, _ := config.Build()
	logger := logInit(cfg.env == "development", nil)

	// Print a confirmation message that the logWriter has started
	logger.Info("Log writer has started successfully.") // logger := zap.Must(zap.NewProduction())
	// if cfg.env == "development" {
	// 	logger = zap.Must(zap.NewDevelopment())
	// }
	dbPool, err := openDB(cfg)
	if cfg.port == "" {
		logger.Info("$PORT has not been set, setting port to :4000")
		cfg.port = "4000"
	}
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.redis.address})
	defer asynqClient.Close()

	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	defer dbPool.Close()
	defer logger.Sync()
	sugar := logger.Sugar()

	//router := chi.NewRouter()

	mailer, err := mailer.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender, asynqClient)
	if err != nil {
		logger.Fatal("Failed to create mailer", zap.Error(err))
	}

	app := &application{
		config:      cfg,
		logger:      logger,
		router:      chi.NewRouter(),
		models:      data.NewModels(dbPool),
		sqlc:        db.New(dbPool),
		mailer:      mailer,
		asynqClient: asynqClient,
	}
	sugar.Info("Database connection estabilished")
	app.router.Use(app.logRequest)
	app.router.Use(app.authenticate)
	app.Routes()
	if err := app.serve(); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
