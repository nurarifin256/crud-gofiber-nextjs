package controllers

import (
	"fmt"
	"path/filepath"
	"simple-crud/configs"
	"simple-crud/helpers"
	"simple-crud/models"
	"simple-crud/requests"
	"simple-crud/utils"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

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

	// debug c.FormFile
	if err != nil {
		return helpers.ResponseJson(c, fiber.StatusBadRequest, "error di sini", err.Error(), []interface{}{})
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

func LoginUser(c *fiber.Ctx) error {
	// init login request
	data := new(requests.LoginRequest)

	// validasi field in body json
	if err := c.BodyParser(data); err != nil {
		return helpers.ResponseJson(c, 400, "warning", err.Error(), []interface{}{})
	}

	// get validator & translator
	validate := helpers.GetValidator()
	translator := helpers.GetTranslator()

	// validate data
	if err := validate.Struct(data); err != nil {
		// parse error validation
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Translate(translator)
		}

		return helpers.ResponseJson(c, 422, "warning", errors, []interface{}{})
	}

	// init db
	db := configs.DB.Db

	// check email
	var user models.User
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return helpers.ResponseJson(c, 401, "warning", "Email not found", []interface{}{})
		}

		return helpers.ResponseJson(c, 500, "warning", err.Error(), []interface{}{})
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return helpers.ResponseJson(c, 401, "warning", "Password incorrect", []interface{}{})
	}

	// check if a token already exists for the user
	var userToken models.UserToken
	if err := db.Where("user_id = ?", user.ID).First(&userToken).Error; err == nil {
		return helpers.ResponseJson(c, 200, "success", "token already exists", userToken.Token)
	}

	// generate token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return helpers.ResponseJson(c, 500, "danger", err.Error(), []interface{}{})
	}

	// mapping data
	userToken = models.UserToken{
		UserID: user.ID,
		Token:  token,
	}

	// save token
	if err := db.Create(&userToken).Error; err != nil {
		return helpers.ResponseJson(c, 500, "danger", err.Error(), []interface{}{})
	}

	// mapping response
	userResponse := fiber.Map{
		"user":  user,
		"token": token,
	}

	return helpers.ResponseJson(c, 200, "success", "Login success", userResponse)
}

func LogoutUser(c *fiber.Ctx) error {
	userId := c.Locals("user").(models.User)
	db := configs.DB.Db
	user := models.UserToken{}

	if err := db.Where("user_id = ?", userId.ID).First(&user).Error; err != nil {
		return helpers.ResponseJson(c, fiber.StatusForbidden, "danger", err.Error(), []interface{}{})
	} else {
		db.Delete(&user)

		return helpers.ResponseJson(c, fiber.StatusOK, "success", "logout success", []interface{}{})
	}
}

func ListUser(c *fiber.Ctx) error {
	var users []UserDTO
	var count int64

	where := "1=1"
	queryArgs := []interface{}{} // query arguments

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	skip, _ := strconv.Atoi(c.Query("skip", "0"))
	search := c.Query("search", "")

	if search != "" {
		where += " AND (username ILIKE ? OR email ILIKE ?)"
		queryArgs = append(queryArgs, "%"+search+"%", "%"+search+"%")
	}

	db := configs.DB.Db
	query := db.Model(&models.User{}).Select("id, username, email")

	if len(queryArgs) > 0 {
		query = query.Where(where, queryArgs...)
	}

	result := query.Limit(limit).Offset(skip).Find(&users)
	db.Model(&models.User{}).Count(&count)

	if result.Error != nil {
		return helpers.ResponseJson(c, fiber.StatusBadRequest, "warning", result.Error.Error(), []interface{}{})
	}

	if result.RowsAffected == 0 {
		return helpers.ResponseJson(c, fiber.StatusOK, "success", "Data not found", []interface{}{})
	}

	response := fiber.Map{
		"users": users,
		"total": count,
	}

	return helpers.ResponseJson(c, fiber.StatusOK, "success", "List user", response)
}
