package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/athifirshad/eucalyptus/internal/tasks"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.config.port,
		Handler:      app.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownError := make(chan error)

	asynqSrv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: app.config.redis.address},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()

    // Register the email task handler
    mux.HandleFunc("email:send", func(ctx context.Context, t *asynq.Task) error {
        var emailTask tasks.EmailTask
        if err := json.Unmarshal(t.Payload(), &emailTask); err != nil {
			app.logger.Error("An error occurred", zap.Error(err))
        }
        // Assuming app is accessible here, or you have another way to access your application's context
        return app.processEmailTask(ctx, t)
    })

    // Start the Asynq server in a goroutine
    go func() {
        if err := asynqSrv.Run(mux); err != nil {
            app.logger.Error("Asynq server failed to start", zap.Error(err))
        }
    }()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.logger.Sugar().Infof("caught signal: %s", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.logger.Sugar().Infof("completing background tasks: addr=%s", srv.Addr)

		app.wg.Wait()
		shutdownError <- nil
	}()
	app.logger.Sugar().Infof("Starting %s server on %s", app.config.env, srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.Sugar().Infof("stopped server: addr=%s", srv.Addr)

	return nil
}
