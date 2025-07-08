package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"demo/internal/handler"
	"demo/internal/repository"
	"demo/internal/service"
)

func main() {
	// Init Fiber
	app := fiber.New()

	// Middlewares
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
