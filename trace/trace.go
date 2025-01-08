// Copyright (c) Loong Zhou.
// All Rights Reserved.

package trace

import (
	"context"
	"os"
	"time"

	"github.com/aloong-planet/opentelemetry-lambda/trace/log"

	lambdadetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const otelTraceEnableEnvVar = "OTEL_TRACE_ENABLED"

var logger *log.ZapLogger

func init() {
	logger = log.NewLogger()
}

func InitTrace(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// Check if tracing is enabled
	if os.Getenv(otelTraceEnableEnvVar) != "true" {
		// Tracing is disabled, return a no-op TracerProvider
		return sdktrace.NewTracerProvider(), nil
	}

	exp, err := otlptracehttp.New(
		ctx,
		// endpoint should be configured with environment variable
		//otlptracehttp.WithEndpoint(os.Getenv("OTLP_TRACE_HTTP_ENDPOINT")),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(time.Millisecond*200),
	)
	if err != nil {
		return nil, err
	}

	detector := lambdadetector.NewResourceDetector()
	res, err := detector.Detect(ctx)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(&RequestSpanProcessor{}),
	)

	// Downstream spans use global tracer provider
	otel.SetTracerProvider(tp)

	// set global propagator to a composite of traceContext and baggage
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{}),
	)

	return tp, nil
}

// GetTracerProvider returns the registered global trace provider.
func GetTracerProvider() trace.TracerProvider {
	return otel.GetTracerProvider()
}

// ShutdownTracerProvider shutdown the tracer provider and force flushes any
// remaining spans in the buffer.
func ShutdownTracerProvider(ctx context.Context, tp *sdktrace.TracerProvider) {
	err := tp.Shutdown(ctx)
	if err != nil {
		logger.Errorf("error shutting down tracer provider: %v", err)
	}
}
