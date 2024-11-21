package gotel

import (
	"context"

	"github.com/rhoat/go-exercise/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, otelConfig config.OtelConfig, options ...CfgOptionFunc) error {
	var shutdownFuncs []func(context.Context) error
	traceExporter, metricExporter, logExporter, err := setupExporters(ctx, otelConfig.Destination)
	if err != nil {
		return err
	}

	cfg, err := NewConfig(traceExporter, metricExporter, logExporter, options...)
	if err != nil {
		return err
	}
	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider := newTraceProvider(*cfg)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider := newMeterProvider(*cfg)
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)
	// Set up log provider.
	loggerProvider := newLoggerProvider(*cfg)

	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	ShutDown = generateShutdownFunc(shutdownFuncs)
	return nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
