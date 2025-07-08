package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"demo/internal/handler"
	"demo/internal/repository"
	"demo/internal/service"
	"demo/pkg/logger"
	"demo/pkg/middleware"
)

func main() {
	// Init Logger
	logger.Init()

	// Init Fiber
	app := fiber.New()

	// Middlewares
	app.Use(middleware.NewObservabilityMiddleware(logger.Default()))
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
