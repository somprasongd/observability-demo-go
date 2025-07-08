package middleware

import (
	"context"
	"demo/pkg/ctxkey"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func NewObservabilityMiddleware(
	baseLogger *zap.Logger,
	tracer trace.Tracer,
) fiber.Handler {

	// Skip Paths ที่ไม่ต้องการ trace
	skipPaths := map[string]bool{
		"/health": true,
	}
	// กรณีมีการ serve SPA
	staticPrefixes := []string{"/static", "/assets", "/public", "/favicon", "/robots.txt"}

	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()

		// ตรวจสอบ path ที่เรียกมา
		skip := skipPaths[path]
		for _, prefix := range staticPrefixes {
			if strings.HasPrefix(path, prefix) {
				skip = true
				break
			}
		}

		var (
			ctx     context.Context
			span    trace.Span
			traceID string
		)

		if skip {
			ctx = c.Context()
		} else {
			ctx, span = tracer.Start(c.Context(), "HTTP "+c.Method()+" "+path)
			defer span.End()
			traceID = span.SpanContext().TraceID().String()
		}

		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Bind Request ID ลง Response Header
		c.Set("X-Request-ID", requestID)

		// สร้าง child logger
		reqLogger := baseLogger.With(
			zap.String("trace_id", traceID), // เพิ่ม trace_id เพื่อเชื่อมโยง log กับ trace
			zap.String("request_id", requestID),
		)

		// สร้าง Context ใหม่ที่มี logger
		ctx = context.WithValue(ctx, ctxkey.Logger{}, reqLogger)
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
