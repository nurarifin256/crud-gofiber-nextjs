package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simple-crud/api/v1"
	"simple-crud/configs"
	"simple-crud/database/migrations"
	logger "simple-crud/pkg/logger"
	"simple-crud/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func main() {
	cfg := configs.Load()
	logr, flush := logger.New(cfg.App.Env)
	defer flush()
	zap.ReplaceGlobals(logr.Desugar())

	// open database connection
	configs.DB = configs.Open(cfg.DB)

	if cfg.DB.Migrate {
		migrations.MigrateAll()
	}

	app := fiber.New(fiber.Config{ErrorHandler: response.ErrorHandler})
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.App.CorsAllowOrigin,
		AllowHeaders: cfg.App.CorsHeaders,
	}))

	// register routes
	api.Init(app)

	// --- Start Fiber (non-blocking) + graceful shutdown
	go func() {
		logr.Infof("listening on %s", cfg.App.Addr())
		if err := app.Listen(cfg.App.Addr()); err != nil {
			logr.Fatalf("app listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// stop http first
	_ = app.Shutdown()

	// give in-flight job up to 3m (runner.Stop() will wait running jobs to finish)
	_, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
}
