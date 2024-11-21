package routes

import "github.com/gin-gonic/gin"

// Route is the interface that all route groups must implement.
type Route interface {
	Register(rg *gin.RouterGroup)
}

// AddRoutes dynamically registers all route groups to the Gin engine.
func AddRoutes(router *gin.Engine, routes []Route) {
	// Registering the API routes under /api
	api := router.Group("/api")

	// Separate group for health routes (readyz and livez)
	health := router.Group("/")

	// Loop through and register routes
	for _, route := range routes {
		// If the route is a health check, register it under the /healthz path
		if _, ok := route.(*HealthRoutes); ok {
			route.Register(health)
		} else {
			// Otherwise, register it under the /api path
			route.Register(api)
		}
	}
}
