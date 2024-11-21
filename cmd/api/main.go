package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-logr/zapr"
	"github.com/rhoat/go-exercise/pkg/config"
	"github.com/rhoat/go-exercise/pkg/routes"
	"github.com/rhoat/go-exercise/pkg/server"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			AsRoute(routes.NewHealthRoutes),
			AsRoute(routes.NewVersionRoutes),
			AsRoute(routes.NewAdminRoutes),
			config.LoadConfig,
			NewLogger,
			fx.Annotate(
				server.New,
				fx.ParamTags(``, ``, ``, `group:"serverRoutes"`),
			),
		),
		fx.Invoke(func(*gin.Engine) {}),
	)
	app.Run()
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(routes.Route)),
		fx.ResultTags(`group:"serverRoutes"`),
	)
}

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(cfg.LogLevel())
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	otel.SetLogger(zapr.NewLogger(logger))
	return logger, nil
}
