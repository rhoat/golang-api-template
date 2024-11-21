package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rhoat/go-exercise/pkg/handlers"
)

type AdminRoutes struct {
	handler *handlers.AdminHandler
}

func NewAdminRoutes() (*AdminRoutes, error) {
	handler, err := handlers.NewAdminHandler()
	if err != nil {
		return nil, err
	}
	return &AdminRoutes{handler: handler}, nil
}

// Register sets up the /api/admin routes.
func (r *AdminRoutes) Register(rg *gin.RouterGroup) {
	rg.GET("/admin", r.handler.AdminFunction)
}
