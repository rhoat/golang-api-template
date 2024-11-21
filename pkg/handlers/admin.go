package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/metric"
)

type AdminHandler struct {
	logger         *slog.Logger
	requestCounter metric.Int64Counter
}

func NewAdminHandler() (*AdminHandler, error) {
	logger := otelslog.NewLogger(pkgName)
	apiCounter, err := meter.Int64Counter(
		"admin.api.counter",
		metric.WithDescription("Number of admin API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		return nil, err
	}
	return &AdminHandler{
		logger:         logger,
		requestCounter: apiCounter,
	}, nil
}

func (h *AdminHandler) AdminFunction(c *gin.Context) {
	ctx := c.Request.Context()
	h.requestCounter.Add(ctx, 1)
	ctx, span := tracer.Start(c.Request.Context(), "AdminFunction")
	defer span.End()
	h.logger.InfoContext(ctx, "Running admin func")
	renderResponse(ctx, c)
}

func renderResponse(ctx context.Context, c *gin.Context) {
	// discarding ctx because we aren't using it here.
	_, childSpan := tracer.Start(ctx, "renderResponse")
	defer childSpan.End()
	c.IndentedJSON(http.StatusOK, gin.H{"adminFunction": "adminFunction content"})
}
