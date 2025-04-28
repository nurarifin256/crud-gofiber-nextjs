package controllers

import (
	"simple-crud/configs"
	"simple-crud/helpers"
	"simple-crud/models"
	"simple-crud/requests"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreatePurchaseOrder(c *fiber.Ctx) error {
	// init db
	db := configs.DB.Db

	// init request struct
	data := new(requests.PurchaseOrderRequest)

	// validate request
	if err := c.BodyParser(data); err != nil {
		return helpers.ResponseJson(c, 400, "warning", err.Error(), []interface{}{})
	}

	// get validator and translator
	validate := helpers.GetValidator()
	translator := helpers.GetTranslator()

	// validate request data
	if err := validate.Struct(data); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Translate(translator)
		}

		return helpers.ResponseJson(c, 422, "warning", errors, []interface{}{})
	}

	// create purchase order instance
	transaction := models.PurchaseOrder{
		CustomerName:    data.CustomerName,
		CustomerAddress: data.CustomerAddress,
		PoNumber:        data.PoNumber,
		TotalAmount:     data.TotalAmount,
		ShippingDate:    data.ShippingDate,
	}

	if err := db.Create(&transaction).Error; err != nil {
		return helpers.ResponseJson(c, 500, "error", err.Error(), []interface{}{})
	}

	response := fiber.Map{
		"transaction": transaction,
	}

	return helpers.ResponseJson(c, 201, "success", "Purchase order created successfully", response)
}
