package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type config struct {
	port string
	env  string
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
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.Parse()

	logger := zap.Must(zap.NewProduction())
	if cfg.env == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	defer logger.Sync()
	sugar := logger.Sugar()

	router := chi.NewRouter()

	app := &application{
		config: cfg,
		logger: logger,
		router: router,
	}
	sugar.Infof("starting %s server on %s", cfg.env, cfg.port)
	app.router.Use(app.logRequest)
	app.Routes()
	http.ListenAndServe(cfg.port, app.router)

}
