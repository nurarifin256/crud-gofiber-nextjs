package response

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// map validator tag -> pesan manusiawi
func humanizeValidation(ve validator.FieldError) string {
	field := ve.Field()
	tag := ve.Tag()
	param := ve.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("The %s field is required.", lower(field))
	case "email":
		return fmt.Sprintf("The %s must be a valid email address.", lower(field))
	case "min":
		// min bisa numeric atau string length; kita fallback "at least N characters"
		if _, err := strconv.Atoi(param); err == nil {
			return fmt.Sprintf("The %s must be at least %s characters.", lower(field), param)
		}
		return fmt.Sprintf("The %s must be at least %s.", lower(field), param)
	case "max":
		if _, err := strconv.Atoi(param); err == nil {
			return fmt.Sprintf("The %s may not be greater than %s characters.", lower(field), param)
		}
		return fmt.Sprintf("The %s may not be greater than %s.", lower(field), param)
	case "uuid", "uuid4":
		return fmt.Sprintf("The %s must be a valid UUID.", lower(field))
	case "numeric": // paling mirip contoh "integer"
		return fmt.Sprintf("The %s must be an integer.", lower(field))
	case "gte":
		return fmt.Sprintf("The %s must be greater than or equal to %s.", lower(field), param)
	case "lte":
		return fmt.Sprintf("The %s must be less than or equal to %s.", lower(field), param)
	case "eqfield":
		return fmt.Sprintf("The %s and %s must be the same.", lower(field), lower(param))
	case "nefield":
		return fmt.Sprintf("The %s and %s may not be the same.", lower(field), lower(param))
	case "eqcsfield":
		return fmt.Sprintf("The %s and %s must be the same.", lower(field), param)
	case "necsfield":
		return fmt.Sprintf("The %s and %s may not be the same.", lower(field), param)
	case "required_with":
		return fmt.Sprintf("The %s field is required when %s is present.", lower(field), lower(param))
	default:
		// fallback bawaan validator (kurang cantik tapi informatif)
		return ve.Error()
	}
}

func lower(s string) string {
	if s == "" {
		return s
	}
	// kecilkan huruf pertama: "Email" -> "email"
	r := []rune(s)
	if r[0] >= 'A' && r[0] <= 'Z' {
		r[0] = r[0] + 32
	}
	return string(r)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	lg := zap.S()
	
	// 1) Fiber HTTP errors -> bentuk sesuai template
	var fe *fiber.Error
	if errors.As(err, &fe) {
		if fe.Code == fiber.StatusUnauthorized {
			lg.Warnw("unauthorized",
				"status", fe.Code, "method", c.Method(), "path", c.Path(), "ip", c.IP(),
			)
			return Unauthorized(c, "Invalid login credentials")
		}
		if fe.Code >= 500 {
			lg.Errorw("fiber http error",
				"status", fe.Code, "msg", fe.Message, "method", c.Method(), "path", c.Path(), "ip", c.IP(),
			)
		} else {
			lg.Warnw("fiber http error",
				"status", fe.Code, "msg", fe.Message, "method", c.Method(), "path", c.Path(), "ip", c.IP(),
			)
		}
		return Error(c, fe.Code, fe.Message)
	}

	// 2) Validator errors -> 422
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		detail := map[string][]string{}
		for _, ve := range verr {
			key := lower(ve.Field())
			detail[key] = append(detail[key], humanizeValidation(ve))
		}
		lg.Warnw("validation failed",
			"method", c.Method(), "path", c.Path(), "ip", c.IP(), "fields", detail,
		)
		return Validation(c, detail)
	}

	// 3) Not found -> 404
	if errors.Is(err, gorm.ErrRecordNotFound) {
		lg.Warnw("record not found",
			"method", c.Method(), "path", c.Path(), "ip", c.IP(),
		)
		return Error(c, fiber.StatusNotFound, "Not Found")
	}

	// 4) Fallback -> 500
	lg.Errorw("unhandled error",
		"error", err, "method", c.Method(), "path", c.Path(), "ip", c.IP(),
	)
	return Internal(c)
}