package main

import (
	"simple-crud/configs"
	"simple-crud/helpers"
	"simple-crud/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	helpers.InitValidator()

	configs.ConnectDb()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")
	routes.SetUpUserRoutes(api)

	app.Listen(":3000")
}
