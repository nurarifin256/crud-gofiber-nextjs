package routes

import (
	"simple-crud/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetUpUserRoutes(group fiber.Router) {
	userRoute := group.Group("/user")

	userRoute.Post("/", controllers.CreateUser)
}
