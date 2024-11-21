package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rhoat/go-exercise/pkg/handlers"
)

type HealthRoutes struct {
	handler *handlers.HealthHandler
}

func NewHealthRoutes() (*HealthRoutes, error) {
	handler, err := handlers.NewHealthHandler()
	if err != nil {
		return nil, err
	}
	return &HealthRoutes{handler: handler}, nil
}

// Register sets up the /api/admin routes.
func (r *HealthRoutes) Register(rg *gin.RouterGroup) {
	rg.GET("/readyz", r.handler.Readyz)
	rg.GET("/livez", r.handler.Livez)
}
