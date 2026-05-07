package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type dailyWriter struct {
	mu       sync.Mutex
	dir      string
	base     string
	currDate string
	f        *os.File
}

func newDailyWriter(dir, base string) (*dailyWriter, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	w := &dailyWriter{dir: dir, base: base}
	if err := w.rotateIfNeeded(time.Now()); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *dailyWriter) filename(t time.Time) (string, string) {
	d := t.Format("20060102")
	name := w.base + "-" + d + ".log" // ex. app-20250828.log
	return filepath.Join(w.dir, name), d
}

func (w *dailyWriter) rotateIfNeeded(now time.Time) error {
	path, d := w.filename(now)
	if w.f != nil && d == w.currDate {
		return nil
	}
	if w.f != nil {
		_ = w.f.Sync()
		_ = w.f.Close()
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	w.f = f
	w.currDate = d
	return nil
}

func (w *dailyWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.rotateIfNeeded(time.Now()); err != nil {
		return 0, err
	}
	return w.f.Write(p)
}

func (w *dailyWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.f != nil {
		err := w.f.Close()
		w.f = nil
		return err
	}
	return nil
}

func New(env string) (*zap.SugaredLogger, func()) {
	logDir := "./logs"
	base := "app"

	// File writer (harian)
	w, err := newDailyWriter(logDir, base)
	if err != nil {
		var l *zap.Logger
		if env == "production" {
			l, _ = zap.NewProduction()
		} else {
			l, _ = zap.NewDevelopment()
		}
		return l.Sugar(), func() { _ = l.Sync() }
	}

	fileEncCfg := zap.NewProductionEncoderConfig()
	fileEncCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEnc := zapcore.NewJSONEncoder(fileEncCfg)

	consoleEncCfg := zap.NewDevelopmentEncoderConfig()
	consoleEncCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEnc := zapcore.NewConsoleEncoder(consoleEncCfg)

	level := zap.InfoLevel
	fileCore := zapcore.NewCore(fileEnc, zapcore.AddSync(w), level)
	consoleCore := zapcore.NewCore(consoleEnc, zapcore.AddSync(os.Stdout), level)

	core := zapcore.NewTee(fileCore, consoleCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger.Sugar(), func() {
		_ = logger.Sync()
		_ = w.Close()
	}
}
