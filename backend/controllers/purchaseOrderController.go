package controllers

import (
	"fmt"
	"net/smtp"
	"os"
	"simple-crud/configs"
	"simple-crud/helpers"
	"simple-crud/models"
	"simple-crud/requests"
	"strings"

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

	// create purchase order items
	var items []models.Item
	for _, reqItem := range data.Items {
		item := models.Item{
			PoId:      transaction.ID,
			ItemName:  reqItem.ItemName,
			ItemCode:  reqItem.ItemCode,
			ItemPrice: reqItem.ItemPrice,
			ItemQty:   reqItem.ItemQty,
		}

		items = append(items, item)
	}

	if len(items) > 0 {
		if err := db.Create(&items).Error; err != nil {
			return helpers.ResponseJson(c, 500, "error", err.Error(), []interface{}{})
		}
	}

	response := fiber.Map{
		"transaction": transaction,
		"items":       items,
	}

	return helpers.ResponseJson(c, 201, "success", "Purchase order created successfully", response)
}

func EmailPurchaseOrders(c *fiber.Ctx) error {
	to := []string{"nurarifin.it@gmail.com"}
	cc := []string{"nur.arifin@adis.co.id"}
	subject := "Purchase Order"
	message := "This is a test email from Fiber and Gorm"

	err := sendEmail(to, cc, subject, message)
	if err != nil {
		return helpers.ResponseJson(c, 500, "error", err.Error(), nil)
	}

	return helpers.ResponseJson(c, 200, "success", "Email sent successfully", nil)
}

func sendEmail(to []string, cc []string, subject, message string) error {
	body := "From : " + os.Getenv("EMAIL_SENDER") + "\n" +
		"To : " + strings.Join(to, ", ") + "\n" +
		"Cc : " + strings.Join(cc, ", ") + "\n" +
		"Subject : " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", os.Getenv("EMAIL_AUTH"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST"))
	smtpAddr := fmt.Sprintf("%s:%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT"))

	err := smtp.SendMail(smtpAddr, auth, os.Getenv("EMAIL_AUTH"), append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
