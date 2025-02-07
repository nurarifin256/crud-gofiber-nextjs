package helpers

import "github.com/gofiber/fiber/v2"

func ResponseJson(c *fiber.Ctx, statusCode int, alert string, message any, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"code":    statusCode,
		"type":    alert,
		"message": message,
		"data":    data,
	})
}
