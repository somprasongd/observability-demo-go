package middleware

import (
	"context"
	"demo/pkg/ctxkey"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewObservabilityMiddleware(baseLogger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()

		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Bind Request ID ลง Response Header
		c.Set("X-Request-ID", requestID)

		// สร้าง child logger
		reqLogger := baseLogger.With(
			zap.String("request_id", requestID),
		)

		// สร้าง Context ใหม่
		ctx := context.WithValue(c.Context(), ctxkey.Logger{}, reqLogger)
		// แทน Context เดิม
		c.SetUserContext(ctx)

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()

		// log unhandle error
		if err != nil {
			reqLogger.Error("an error occurred",
				zap.Any("error", err),
				zap.ByteString("stack", debug.Stack()),
			)
		}

		msg := fmt.Sprintf("%d - %s %s", status, method, path)
		reqLogger.Info(msg,
			zap.Int("status", status),
			zap.Float64("duration_sec", duration),
		)

		return err
	}
}
