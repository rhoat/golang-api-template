package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rhoat/go-exercise/pkg/version"
	"go.uber.org/zap"
)

type VersionHandler struct {
	log *zap.Logger
}

func NewVersionHandler(log *zap.Logger) *VersionHandler {
	return &VersionHandler{log: log}
}

func (h *VersionHandler) VersionFunction(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, version.NewInfo())
}
