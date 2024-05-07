package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/athifirshad/eucalyptus/internal/tasks"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readIDParam(r *http.Request) (int32, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return int32(id), nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}

func (app *application) background(fn func()) {
	// Launch a background goroutine.
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.Sugar().Error(fmt.Errorf("%s", err), nil)
			}
		}()
		fn()
	}()
}

func (app *application) DocsHandler(w http.ResponseWriter, r *http.Request) {
	docs := docgen.MarkdownRoutesDoc(app.router, docgen.MarkdownOpts{
		ProjectPath: "github.com/athifirshad/eucalyptus",
		// Intro text included at the top of the generated markdown file.
		Intro: "Generated documentation for Eucalyptus",
	})
	err := os.WriteFile("./docs/routes.md", []byte(docs), 0644)
	if err != nil {
		app.logger.Sugar().Errorf("error writing docs", "err", err)
		http.Error(w, "Internal server error fetching docs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Documentation successfully written to routes.md"))
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

func ReadQr(imagePath string) (string, error) {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Create a BinaryBitmap from the image
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("error creating BinaryBitmap: %v", err)
	}

	// Create a QRCodeReader
	reader := qrcode.NewQRCodeReader()

	// Decode the QR code
	result, err := reader.Decode(bitmap, nil)
	if err != nil {
		return "", fmt.Errorf("error decoding QR code: %v", err)
	}

	// Return the decoded text
	return result.GetText(), nil
}