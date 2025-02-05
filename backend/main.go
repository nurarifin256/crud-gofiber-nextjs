package main

import (
	"simple-crud/configs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.ConnectDb()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
