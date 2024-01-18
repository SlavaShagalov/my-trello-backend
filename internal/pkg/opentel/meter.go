package opentel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	exp, err := NewMetricOTLPExporter(context.Background())
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exp, metric.WithInterval(time.Second))),
	)
	return meterProvider, nil
}

func NewMetricOTLPExporter(ctx context.Context) (metric.Exporter, error) {
	insecureOpt := otlpmetrichttp.WithInsecure()
	endpointOpt := otlpmetrichttp.WithEndpoint(otlpEndpoint)

	return otlpmetrichttp.New(ctx, insecureOpt, endpointOpt)
}
