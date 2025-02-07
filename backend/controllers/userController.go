package controllers

import (
	"fmt"
	// "simple-crud/configs"
	"simple-crud/helpers"
	"simple-crud/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	// db := configs.DB.Db
	user := models.User{}

	user.Username = c.FormValue("username")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")

	file, err := c.FormFile("image")
	if err != nil {
		return helpers.ResponseJson(c, fiber.StatusBadRequest, "Error", err.Error(), []interface{}{})
	}

	currentDate := time.Now().Format("20060102")
	fileName := fmt.Sprintf("%s-%s", currentDate, file.Filename)

	user.Image = fileName

	validate := helpers.GetValidator()
	translator := helpers.GetTranslator()

	if err := validate.Struct(user); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Translate(translator)
		}

		return helpers.ResponseJson(c, 422, "warning", errors, []interface{}{})
	}

	return helpers.ResponseJson(c, 200, "success", "User created successfully", []interface{}{})
}
