package main

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"

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

	// เก็บ Process Metrics: สร้าง Runtime Instrument → ผูกกับ MeterProvider
	runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second * 10),
	)

	// Init Fiber
	app := fiber.New()

	// Middlewares
	app.Use(middleware.NewObservabilityMiddleware(
		logger.Default(),
		otel.TracerProvider.Tracer("demo-app"),
		otel.MeterProvider.Meter("demo-app"), // เพิ่มส่ง meter
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
