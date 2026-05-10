package helpers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// V is the global validator instance
var V = validator.New()

// ValidationConfig holds the global validation configuration
type ValidationConfig struct {
	FieldChecker     FieldUniqueChecker
	CompositeChecker CompositeUniqueChecker
	Context          *ValidationContext
}

// ValidationContext holds context information for validation
type ValidationContext struct {
	ExcludeID *int64
}

// FieldUniqueChecker defines interface for single field unique validation
type FieldUniqueChecker interface {
	FieldExists(ctx context.Context, fieldName, value string, excludeID ...int64) (bool, error)
}

// CompositeUniqueChecker defines interface for composite field unique validation
type CompositeUniqueChecker interface {
	CompositeExists(ctx context.Context, fields map[string]string, excludeID ...int64) (bool, error)
}

// Global validation configuration
var globalValidationConfig = &ValidationConfig{}

// SetFieldChecker sets the field unique checker for validation
func SetFieldChecker(checker FieldUniqueChecker) {
	globalValidationConfig.FieldChecker = checker
}

// SetCompositeChecker sets the composite unique checker for validation
func SetCompositeChecker(checker CompositeUniqueChecker) {
	globalValidationConfig.CompositeChecker = checker
}

// SetValidationContext sets the global validation context
func SetValidationContext(ctx *ValidationContext) {
	globalValidationConfig.Context = ctx
}

// universalUniqueValidator is a universal validator for unique validation
// Supports two modes:
// - unique=field validates single field using FieldExists
// - unique=field1:field2:field3 validates composite fields using CompositeExists
func universalUniqueValidator(fl validator.FieldLevel) bool {
	param := fl.Param()
	if param == "" {
		return false // Parameter is required for unique validation
	}

	// Parse and clean field names - use colon separator only
	fieldNames := strings.Split(param, ":")
	for i := range fieldNames {
		fieldNames[i] = strings.TrimSpace(fieldNames[i])
	}

	// Prepare excludeID for update operations
	var excludeID []int64
	if globalValidationConfig.Context != nil && globalValidationConfig.Context.ExcludeID != nil {
		excludeID = []int64{*globalValidationConfig.Context.ExcludeID}
	}

	// Single field validation mode
	if len(fieldNames) == 1 {
		return validateSingleField(fl, fieldNames[0], excludeID)
	}

	// Composite fields validation mode
	return validateCompositeFields(fl, fieldNames, excludeID)
}

// validateSingleField handles single field unique validation
func validateSingleField(fl validator.FieldLevel, fieldName string, excludeID []int64) bool {
	if globalValidationConfig.FieldChecker == nil {
		return true // Skip validation if no field checker is set
	}

	value := strings.TrimSpace(fl.Field().String())
	if value == "" {
		return true // Skip validation for empty values
	}

	exists, err := globalValidationConfig.FieldChecker.FieldExists(context.Background(), fieldName, value, excludeID...)
	return err == nil && !exists // Valid if no error and not exists
}

// validateCompositeFields handles composite fields unique validation
func validateCompositeFields(fl validator.FieldLevel, fieldNames []string, excludeID []int64) bool {
	if globalValidationConfig.CompositeChecker == nil {
		return true // Skip validation if no composite checker is set
	}

	fieldMap := make(map[string]string, len(fieldNames))
	structValue := fl.Top()

	// Handle pointer dereferencing if needed
	if structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			return true // Skip validation if struct is nil
		}
		structValue = structValue.Elem()
	}

	// Extract field values using database column names directly
	for _, dbColumnName := range fieldNames {
		// Find struct field by matching json tag
		structFieldName := findStructFieldByTag(structValue.Type(), "json", dbColumnName)
		if structFieldName == "" {
			return false // Field not found in struct
		}

		fieldValue := structValue.FieldByName(structFieldName)
		if !fieldValue.IsValid() {
			return false // Field not found in struct
		}

		strValue := extractStringValue(fieldValue)
		if strValue == "" {
			return true // Skip validation if any field is empty
		}

		// Use database column name directly in the map
		fieldMap[dbColumnName] = strValue
	}

	exists, err := globalValidationConfig.CompositeChecker.CompositeExists(context.Background(), fieldMap, excludeID...)
	return err == nil && !exists // Valid if no error and not exists
}

// findStructFieldByTag finds struct field name by matching tag value
func findStructFieldByTag(structType reflect.Type, tagName, tagValue string) string {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if tag := field.Tag.Get(tagName); strings.Split(tag, ",")[0] == tagValue {
			return field.Name
		}
	}
	return ""
}

// extractStringValue converts reflect.Value to string
func extractStringValue(fieldValue reflect.Value) string {
	switch fieldValue.Kind() {
	case reflect.String:
		return strings.TrimSpace(fieldValue.String())
	case reflect.Ptr:
		if !fieldValue.IsNil() {
			elem := fieldValue.Elem()
			switch elem.Kind() {
			case reflect.String:
				return strings.TrimSpace(elem.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return strings.TrimSpace(fmt.Sprintf("%d", elem.Int()))
			default:
				return strings.TrimSpace(fmt.Sprintf("%v", elem.Interface()))
			}
		}
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strings.TrimSpace(fmt.Sprintf("%d", fieldValue.Int()))
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", fieldValue.Interface()))
	}
}

type ValidationErrorResponse struct {
	Code    int                 `json:"code"`
	Type    string              `json:"type"`
	Message map[string][]string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewValidationError(field string, msg string) ValidationErrorResponse {
	return ValidationErrorResponse{
		Code: 422,
		Type: "warning",
		Message: map[string][]string{
			field: {msg},
		},
	}
}

func NewError(msg string) ErrorResponse {
	return ErrorResponse{
		Code:    422,
		Type:    "warning",
		Message: msg,
	}
}

// FormatValidationErrors formats validation errors into a readable format
func FormatValidationErrors(err error) map[string][]string {
	errors := make(map[string][]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())
		var message string

		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("The %s field is required.", field)
		case "min":
			message = fmt.Sprintf("The %s must be at least %s characters.", field, e.Param())
		case "max":
			message = fmt.Sprintf("The %s may not be greater than %s characters.", field, e.Param())
		case "email":
			message = fmt.Sprintf("The %s must be a valid email address.", field)
		case "unique":
			param := e.Param()
			if strings.Contains(param, ":") {
				// Composite unique validation - show combination message
				message = "The combination must be unique."
			} else {
				// Single field unique validation
				message = fmt.Sprintf("%s has already been taken.", e.Field())
			}
		default:
			message = fmt.Sprintf("The %s field is invalid.", field)
		}

		errors[field] = append(errors[field], message)
	}

	return errors
}

// init registers the custom unique validator
func init() {
	V.RegisterValidation("unique", universalUniqueValidator)
}
