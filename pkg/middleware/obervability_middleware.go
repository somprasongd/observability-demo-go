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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func NewObservabilityMiddleware(
	baseLogger *zap.Logger,
	tracer trace.Tracer,
	meter metric.Meter,
) fiber.Handler {

	// ----- OTel Instruments -----
	requestCounter, _ := meter.Int64Counter("http_requests_total")
	requestDuration, _ := meter.Float64Histogram("http_request_duration_ms")
	inflightCounter, _ := meter.Int64UpDownCounter("http_requests_inflight")
	requestSize, _ := meter.Float64Histogram("http_request_size_bytes")
	responseSize, _ := meter.Float64Histogram("http_response_size_bytes")
	errorCounter, _ := meter.Int64Counter("http_requests_error_total")

	// Skip Paths ที่ไม่ต้องการ trace
	skipPaths := map[string]bool{
		"/health":  true,
		"/metrics": true,
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

		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Bind Request ID ลง Response Header
		c.Set("X-Request-ID", requestID)

		var (
			ctx     context.Context
			span    trace.Span
			traceID string
			spanID  string
		)

		if skip {
			ctx = c.Context()
		} else {
			ctx, span = tracer.Start(c.Context(), "HTTP "+c.Method()+" "+path,
				trace.WithAttributes(attribute.String("http.request_id", requestID)),
				trace.WithAttributes(attribute.String("http.request.method", method)),
				trace.WithAttributes(attribute.String("url.path", path)),
			)
			defer span.End()
			traceID = span.SpanContext().TraceID().String()
			spanID = span.SpanContext().SpanID().String()
		}

		// สร้าง child logger
		reqLogger := baseLogger.With(
			zap.String("request_id", requestID),
			zap.String("http.request.method", method),
			zap.String("url.path", path),
		)

		// สร้าง Context ใหม่ที่มี logger
		ctx = context.WithValue(ctx, ctxkey.Logger{}, reqLogger)
		// แทน Context เดิม
		c.SetUserContext(ctx)

		// ----- Record Inflight -----
		if !skip {
			inflightCounter.Add(ctx, 1)
		}

		err := c.Next()

		duration := time.Since(start).Milliseconds()
		status := c.Response().StatusCode()

		if !skip {
			labels := []attribute.KeyValue{
				attribute.String("http.request.method", method),
				attribute.String("url.path", path),
				attribute.Int("http.response.status_code", status),
			}

			requestCounter.Add(ctx, 1, metric.WithAttributes(labels...))
			requestDuration.Record(ctx, float64(duration), metric.WithAttributes(labels...))
			inflightCounter.Add(ctx, -1)

			// Request Size (Header Content-Length)
			if reqSize := c.Request().Header.ContentLength(); reqSize > 0 {
				requestSize.Record(ctx, float64(reqSize), metric.WithAttributes(labels...))
			}

			// Response Size (Body Length)
			if resSize := len(c.Response().Body()); resSize > 0 {
				responseSize.Record(ctx, float64(resSize), metric.WithAttributes(labels...))
			}

			if status >= 400 {
				errorCounter.Add(ctx, 1, metric.WithAttributes(labels...))
			}
		}

		// log unhandle error
		if err != nil {
			reqLogger.Error("an error occurred",
				zap.Any("error", err),
				zap.ByteString("stack", debug.Stack()),
				zap.String("trace_id", traceID), // เพิ่ม trace_id เพื่อเชื่อมโยง log กับ trace
				zap.String("span_id", spanID),   // เพิ่ม span_id เพื่อเชื่อมโยง log กับ trace
			)
		}

		msg := fmt.Sprintf("%d - %s %s", status, method, path)
		reqLogger.Info(msg,
			zap.Int("http.response.status_code", status),
			zap.Int64("duration_ms", duration),
			zap.String("trace_id", traceID), // เพิ่ม trace_id เพื่อเชื่อมโยง log กับ trace
			zap.String("span_id", spanID),   // เพิ่ม span_id เพื่อเชื่อมโยง log กับ trace
		)

		span.SetAttributes(
			attribute.Int("http.response.status_code", status),
		)

		if status >= 400 {
			span.SetStatus(codes.Error, "")
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return err
	}
}
