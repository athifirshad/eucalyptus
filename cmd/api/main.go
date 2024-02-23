package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/athifirshad/eucalyptus/db"
	"github.com/athifirshad/eucalyptus/internal/data"
	"github.com/athifirshad/eucalyptus/internal/mailer"
	"github.com/athifirshad/eucalyptus/internal/tasks"
	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natefinch/lumberjack"
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

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("| " + r.Method + " | " + r.URL.String() + " | " + r.RemoteAddr + " | " + r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

func logInit(d bool, f *os.File) *zap.Logger {
	pe := zap.NewProductionEncoderConfig()

	// Set up lumberjack for log rotation
	currentDate := time.Now().Format("02-01-2006")

	logWriter := &lumberjack.Logger{
		Filename:   "logs/" + currentDate + ".log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	pe.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	fileEncoder := zapcore.NewJSONEncoder(pe)

	level := zap.InfoLevel
	if d {
		level = zap.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logWriter), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	l := zap.New(core) // Creating the logger

	return l
}
func (app *application) processEmailTask(ctx context.Context, task *asynq.Task) error {
	var t tasks.EmailTask
	if err := json.Unmarshal(task.Payload(), &t); err != nil {
		return err
	}
	return app.mailer.SendEmail(t.Recipient, t.Template, t.Data)
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
	flag.StringVar(&cfg.redis.address, "redis-address", "localhost:6379", "Redis address")
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

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.redis.address})
	defer asynqClient.Close()

	
	if err != nil {
		logger.Fatal("Failed to open DB", zap.Error(err))
	}

	defer dbPool.Close()
	defer logger.Sync()
	sugar := logger.Sugar()

	//router := chi.NewRouter()

	mailer, err := mailer.NewMailer(cfg.smtp.host,cfg.smtp.port,cfg.smtp.username,cfg.smtp.password,cfg.smtp.sender, asynqClient)
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
