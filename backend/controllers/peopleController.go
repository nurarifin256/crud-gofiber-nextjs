package controllers

import (
	"mime/multipart"
	"simple-crud/helpers"

	"github.com/gofiber/fiber/v2"
)

const chunkSize = 10000

func getFile(c *fiber.Ctx) (multipart.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	return src, nil
}

func CreateItem(c *fiber.Ctx) error {
	// var wg sync.WaitGroup

	file, err := getFile(c)
	if err != nil {
		return helpers.ResponseJson(c, 400, "warning", err.Error(), []interface{}{})
	}
	defer file.Close()
}
