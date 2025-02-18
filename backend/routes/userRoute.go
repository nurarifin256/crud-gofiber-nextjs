package routes

import (
	"simple-crud/controllers"
	"simple-crud/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpUserRoutes(group fiber.Router) {
	userRoute := group.Group("/user")

	userRoute.Post("/", controllers.CreateUser)
	userRoute.Post("/login", controllers.LoginUser)
	userRoute.Post("/logout", middleware.JWTProtected(), controllers.LogoutUser)
}
