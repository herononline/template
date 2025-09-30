package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-gorote/gorote"
	"github.com/go-gorote/template/env"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/swagger"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	_ "github.com/go-gorote/template/docs"
)

func gracefulShutdown(ctx context.Context, app *fiber.App) {
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
}

// @title           Auth API
// @version         1.0
// @description     API de autenticação
// @termsOfService  http://swagger.io/terms/

// @BasePath  /api/v1

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := fiber.New(fiber.Config{
		AppName:      env.App.Name,
		ServerHeader: env.App.Name,
		ErrorHandler: gorote.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     env.CORS,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.Use(helmet.New(helmet.Config{
		XSSProtection: "1; mode=block",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(env.App.Name),
			semconv.ServiceVersionKey.String(env.App.Version),
		),
	)
	if err != nil {
		log.Fatal("error on resource")
	}

	if err := gorote.TelemetryFiber(ctx, app, res, env.CollectorOpenTelemetry); err != nil {
		log.Fatal("error on telemetry")
	}

	if err := Config(ctx, app); err != nil {
		log.Fatal("error on config")
	}

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", env.App.Port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go gracefulShutdown(ctx, app)

	<-ctx.Done()
	log.Println("Graceful shutdown complete.")
}
