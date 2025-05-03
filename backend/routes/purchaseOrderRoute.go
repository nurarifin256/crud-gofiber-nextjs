package routes

import (
	"simple-crud/controllers"
	"simple-crud/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpPurchaseOrderRoutes(group fiber.Router) {
	purchaseOrderRoute := group.Group("/purchase-order")

	purchaseOrderRoute.Post("/", middleware.JWTProtected(), controllers.CreatePurchaseOrder)
	purchaseOrderRoute.Get("/email", middleware.JWTProtected(), controllers.EmailPurchaseOrders)
}
