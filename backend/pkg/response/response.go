package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Body struct {
	Code    int         `json:"code"`
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

// --- penentuan type berdasar status ---
func kindOf(code int) string {
	switch {
	case code >= 500:
		return "danger"
	case code >= 400:
		return "warning"
	case code >= 200:
		return "success"
	default:
		return "normal"
	}
}

func JSON(c *fiber.Ctx, code int, message interface{}, data interface{}, typeOpsional ...string) error {
	responseType := kindOf(code) // default

	if len(typeOpsional) > 0 && typeOpsional[0] != "" {
		responseType = typeOpsional[0]
	}

	b := Body{
		Code:    code,
		Type:    responseType,
		Message: message,
		Data:    data,
	}
	return c.Status(code).JSON(b)
}

// --- Success (200/201) ---
func OK(c *fiber.Ctx, msg string, data interface{}) error {
	return JSON(c, http.StatusOK, msg, data)
}

func Created(c *fiber.Ctx, data interface{}) error {
	return JSON(c, http.StatusCreated, "Successfully Created!", data)
}

// --- 401 Unauthorized ---
func Unauthorized(c *fiber.Ctx, msg ...string) error {
	m := "Invalid login credentials"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}
	return JSON(c, http.StatusUnauthorized, m, []any{})
}

// --- 422 Unprocessable Entity ---
func Validation(c *fiber.Ctx, fieldErrors map[string][]string) error {
	return JSON(c, http.StatusUnprocessableEntity, fieldErrors, []any{})
}

// --- 500 Internal Server Error ---
func Internal(c *fiber.Ctx, msg ...string) error {
	m := "Oops, an error has accurred, we'll be up and running shortly. if you need immediate help, please call us"
	if len(msg) > 0 && msg[0] != "" {
		m = msg[0]
	}
	return JSON(c, http.StatusInternalServerError, m, []any{})
}

// --- Helper untuk error bebas kode ---
func Error(c *fiber.Ctx, code int, msg interface{}) error {
	// untuk error, data selalu array kosong sesuai template
	return JSON(c, code, msg, []any{})
}

func UnauthorizedErr(msg string) error {
	if msg == "" {
		msg = "Incorrect NIK or Password"
	}
	return fiber.NewError(fiber.StatusUnauthorized, msg)
}

func Errorf(status int, msg string) error {
	if msg == "" {
		msg = http.StatusText(status)
	}
	return fiber.NewError(status, msg)
}
