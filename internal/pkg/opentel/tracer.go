package opentel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	otlpEndpoint string
)

func init() {
	otlpEndpoint = "jaeger:4318"
	//otlpEndpoint = "127.0.0.1:4318"
}

func NewTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	//exp, err := NewConsoleExporter()
	exp, err := NewOTLPExporter(context.Background())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func NewConsoleExporter() (trace.SpanExporter, error) {
	return stdouttrace.New()
}

func NewOTLPExporter(ctx context.Context) (trace.SpanExporter, error) {
	insecureOpt := otlptracehttp.WithInsecure()
	endpointOpt := otlptracehttp.WithEndpoint(otlpEndpoint)

	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}
