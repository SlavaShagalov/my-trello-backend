package opentel

import (
	"context"
	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

var (
	Tracer  trace.Tracer
	Meter   metric.Meter
	Counter metric.Int64Counter
)

func NewResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func SetupOTelSDK(ctx context.Context, logger *zap.Logger, serviceName, serviceVersion string) (*sdktrace.TracerProvider,
	*sdkmetric.MeterProvider, error) {
	res, err := NewResource(serviceName, serviceVersion)
	if err != nil {
		logger.Error("Failed to initialize resource", zap.Error(err))
		return nil, nil, err
	}

	tp, err := NewTraceProvider(res)
	if err != nil {
		logger.Error("Failed to create TraceProvider", zap.Error(err))
		return nil, nil, err
	}
	otel.SetTracerProvider(tp)

	mp, err := NewMeterProvider(res)
	if err != nil {
		logger.Error("Failed to create MeterProvider", zap.Error(err))
		return nil, nil, err
	}
	otel.SetMeterProvider(mp)

	Tracer = tp.Tracer("tracer")
	Meter = mp.Meter("meter")

	Counter, err = Meter.Int64Counter("my.Counter",
		metric.WithDescription("My Counter description"),
		metric.WithUnit("{Counter}"))
	if err != nil {
		return nil, nil, err
	}

	logger.Info("OpenTelemetry setup")
	return tp, mp, nil
}
