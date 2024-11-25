package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rhoat/go-exercise/pkg/config"
	"github.com/rhoat/go-exercise/pkg/middleware"
	gotel "github.com/rhoat/go-exercise/pkg/otel"
	"github.com/rhoat/go-exercise/pkg/routes"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, log *zap.Logger, cfg *config.Config, routeGroups []routes.Route) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // Optional: suppress Gin debug output
	router := gin.New()
	router.Use(middleware.Logging(log))
	router.Use(gin.Recovery())

	// Create route instances
	routes.AddRoutes(router, routeGroups)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerConfig.Port,
		Handler:           router,
		ReadTimeout:       time.Duration(cfg.ServerConfig.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.ServerConfig.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.ServerConfig.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.ServerConfig.IdleTimeout) * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("Setting up Otel SDK", zap.String("destination", cfg.Otel.Destination.String()))
			err := gotel.SetupOTelSDK(context.Background(), cfg.Otel.Destination)
			if err != nil {
				log.Error("Failed to setup otel", zap.Error(err))
				return err
			}
			ln, err := net.Listen("tcp", srv.Addr) // the web server starts listening on 8080
			if err != nil {
				log.Error("Failed to start HTTP Server", zap.String("Addr", srv.Addr), zap.Error(err))
				return err
			}
			go func() {
				if err = srv.Serve(ln); err != nil {
					log.Error("HTTP Server error", zap.Error(err))
				}
			}()
			log.Info("started HTTP Server", zap.String("Addr", srv.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := gotel.ShutDown(ctx); err != nil {
				log.Error("shutting down otel", zap.Error(err))
			}
			if err := srv.Shutdown(ctx); err != nil {
				log.Error("shutting down http server", zap.String("Addr", srv.Addr), zap.Error(err))
			}
			log.Info("stopped HTTP Server", zap.String("Addr", srv.Addr))
			return nil
		},
	})

	return router
}
