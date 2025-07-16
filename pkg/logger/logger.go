package logger

import (
	"context"
	"demo/pkg/ctxkey"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var baseLogger *zap.Logger

func Init() {
	config := zap.NewProductionConfig()
	// ตั้งค่า key ของ field ให้ตรงกับ OTel semantic convention
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "severity"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	l, err := config.Build()

	if err != nil {
		panic(err)
	}

	baseLogger = l.With(zap.String("service.name", "demo-app"))
}

func Default() *zap.Logger {
	return baseLogger
}

// FromContext extracts logger from context
func FromContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(ctxkey.Logger{}).(*zap.Logger)
	if !ok {
		log = baseLogger
	}

	// ดึง trace_id + span_id จาก OTel context
	span := trace.SpanContextFromContext(ctx)

	if span.IsValid() {
		return log.With(
			zap.String("trace_id", span.TraceID().String()), // เพิ่ม trace_id เพื่อเชื่อมโยง log กับ trace
			zap.String("span_id", span.SpanID().String()),   // เพิ่ม span_id เพื่อเชื่อมโยง log กับ trace
		)
	}

	return log
}
