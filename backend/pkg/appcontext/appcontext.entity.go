package appcontext

import (
	"github.com/gofiber/fiber/v2"
)

type RouteDefinition struct {
	Method         string
	Path           string
	Action         fiber.Handler
	Middleware     []fiber.Handler
	AuthRequired   bool
	AuthMiddleware []fiber.Handler
}

// RouteInitializer interface defines the contract for route initializers
type RouteInitializer interface {
	Initialize() []RouteDefinition
}
