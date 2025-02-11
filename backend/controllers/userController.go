package controllers

import (
	"fmt"
	"path/filepath"
	"simple-crud/configs"
	"simple-crud/helpers"
	"simple-crud/models"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	// init db
	db := configs.DB.Db

	// init user model
	user := models.User{}

	// get form value
	user.Username = c.FormValue("username")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")

	file, err := c.FormFile("image")
	if err != nil {
		return helpers.ResponseJson(c, fiber.StatusBadRequest, "Error", err.Error(), []interface{}{})
	}

	// create format filename save to database
	currentDate := time.Now().Format("20060102")
	fileName := fmt.Sprintf("%s-%s", currentDate, file.Filename)

	user.Image = fileName

	// init validator
	validate := helpers.GetValidator()
	translator := helpers.GetTranslator()

	// validate user
	if err := validate.Struct(user); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Translate(translator)
		}

		return helpers.ResponseJson(c, 422, "warning", errors, []interface{}{})
	}

	// encrypt password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.ResponseJson(c, 400, "warning", err.Error(), []interface{}{})
	}
	user.Password = string(hashPassword)

	if err := db.Create(&user).Error; err != nil {
		return helpers.ResponseJson(c, 400, "warning", err.Error(), []interface{}{})
	}

	// clean the filename
	fileName = strings.ReplaceAll(fileName, " ", "_")
	fileName = filepath.Base(fileName)

	// save file to folder images
	savePath := fmt.Sprintf("./images/%s", fileName)
	if err := c.SaveFile(file, savePath); err != nil {
		return helpers.ResponseJson(c, 400, "warning", err.Error(), []interface{}{})
	}

	return helpers.ResponseJson(c, 200, "success", "User created successfully", user)
}
