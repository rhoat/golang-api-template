package handlers

import (
	"go.opentelemetry.io/otel"
)

var (
	pkgName = "github.com/rhoat/go-exercise/pkg/handlers"
	meter   = otel.Meter(pkgName)
	tracer  = otel.Tracer(pkgName)
)
