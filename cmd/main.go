package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"demo/internal/handler"
	"demo/internal/repository"
	"demo/internal/service"
	"demo/pkg/logger"
	"demo/pkg/middleware"
	"demo/pkg/observability"
)

func main() {
	// Init Logger
	logger.Init()

	// Init Observability via Opentelmetry
	otel, err := observability.NewOTel(
		context.Background(),
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		"demo-app")
	if err != nil {
		logger.Default().Fatal(err.Error())
	}
	defer otel.Shutdown(context.Background())

	// Init Fiber
	app := fiber.New()

	// Middlewares
	app.Use(middleware.NewObservabilityMiddleware(
		logger.Default(),
		otel.TracerProvider.Tracer("demo-app"), // เพิ่มส่ง tracer
	))
	app.Use(cors.New())
	app.Use(recover.New())

	// Init DI
	repo := repository.NewUserRepository()
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// Routes
	app.Get("/users/:id", h.GetUser)

	// Start
	app.Listen(":8080")
}
