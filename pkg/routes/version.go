package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rhoat/go-exercise/pkg/handlers"
	"go.uber.org/zap"
)

type VersionRoutes struct {
	handler *handlers.VersionHandler
}

func NewVersionRoutes(log *zap.Logger) (*VersionRoutes, error) {
	return &VersionRoutes{handler: handlers.NewVersionHandler(log)}, nil
}

// Register sets up the /api/admin routes.
func (r *VersionRoutes) Register(rg *gin.RouterGroup) {
	admin := rg.Group("/version")
	admin.GET("/", r.handler.VersionFunction)
}
