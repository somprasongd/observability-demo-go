package handler

import (
	"demo/internal/service"
	"demo/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ctx := c.UserContext()
	logger := logger.FromContext(ctx)

	tracer := trace.SpanFromContext(ctx).TracerProvider().Tracer("handler")
	ctx, span := tracer.Start(ctx, "Handler:GetUser")
	defer span.End()

	logger.Info("Handler: GetUser called")

	id := c.Params("id")

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		logger.Error("Handler: Failed to get user", zap.Error(err))
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(user)
}
