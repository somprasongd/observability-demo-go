package handler

import (
	"demo/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	log.Println("Handler: GetUser called")

	id := c.Params("id")

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		log.Println("Handler: Failed to get user")
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(user)
}
