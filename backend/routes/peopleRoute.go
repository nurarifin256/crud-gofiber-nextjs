package routes

import (
	"simple-crud/controllers"
	"simple-crud/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpPeopleRoutes(group fiber.Router) {
	peopleRoute := group.Group("/people")

	peopleRoute.Post("/", middleware.JWTProtected(), controllers.CreateItem)
}
