package api

import (
	"simple-crud/api/v1/routes"

	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	// Initialize API version 1 routes
	routes.RouteV1(app)
}
