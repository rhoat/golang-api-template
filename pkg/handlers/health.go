package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rhoat/go-exercise/pkg/health"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/metric"

	"github.com/gin-gonic/gin"
)

type route string

const (
	livez  route = "livez"
	readyz route = "readyz"
)

type HealthHandler struct {
	logger         *slog.Logger
	checks         *health.Registry
	requestCounter map[route]metric.Int64Counter
}

func initializeMestricCounters() (map[route]metric.Int64Counter, error) {
	result := make(map[route]metric.Int64Counter)
	for _, route := range []route{livez, readyz} {
		counter, err := meter.Int64Counter(
			fmt.Sprintf("health.%s.api.counter", route),
			metric.WithDescription(fmt.Sprintf("Number of %s API calls.", route)),
			metric.WithUnit("{call}"),
		)
		if err != nil {
			return nil, err
		}
		result[route] = counter
	}
	return result, nil
}

func NewHealthHandler() (*HealthHandler, error) {
	logger := otelslog.NewLogger(pkgName)
	counters, err := initializeMestricCounters()
	if err != nil {
		return nil, err
	}

	return &HealthHandler{
		requestCounter: counters,
		logger:         logger,
		checks:         health.NewRegistry(),
	}, nil
}

func (h *HealthHandler) Livez(c *gin.Context) {
	ctx := c.Request.Context()
	h.requestCounter[livez].Add(ctx, 1)
	ctx, span := tracer.Start(c.Request.Context(), "Livez")
	defer span.End()
	h.logger.InfoContext(ctx, "Running Livez checks")
	results := h.checks.RunChecks(ctx)
	statusCode := http.StatusOK

	for _, result := range results {
		if result.Status == health.Fail {
			statusCode = http.StatusInternalServerError
			break
		}
	}

	c.IndentedJSON(statusCode, gin.H{"status": "Alive"})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	ctx := c.Request.Context()
	h.requestCounter[readyz].Add(ctx, 1)
	ctx, span := tracer.Start(c.Request.Context(), "Readyz")
	defer span.End()
	h.logger.InfoContext(ctx, "Running Readyz checks")
	exclude := c.Query("exclude")

	results := h.checks.RunChecks(ctx)
	filteredResults := []health.Status{}

	for _, result := range results {
		if result.Name == exclude {
			continue
		}
		filteredResults = append(filteredResults, result)
	}

	statusCode := http.StatusOK
	for _, result := range filteredResults {
		if result.Status == health.Fail {
			statusCode = http.StatusInternalServerError
			break
		}
	}

	c.IndentedJSON(statusCode, filteredResults)
}
