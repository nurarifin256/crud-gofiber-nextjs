package routes

import (
	v1_auth "simple-crud/api/v1/routes/authentications"
	"simple-crud/pkg/appcontext"

	"github.com/gofiber/fiber/v2"
)

// InitializeRoutes initializes all API routes for version 1
func RouteV1(app *fiber.App) {
	// Initialize route initializer
	initializer := appcontext.NewRouteInitializer()
	// Initialize all route modules
	initializer.InitializeAllRoutes()
	// Register routes to Fiber app
	appcontext.RegisterRoutesToApp(app, "/api/v1")
}

// init function to register route initializers
func init() {
	// Register authentication routes
	appcontext.RegisterRouteInitializer("auth", v1_auth.NewAuthRouteInitializer())
}
