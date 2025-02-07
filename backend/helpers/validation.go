package helpers

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func InitValidator() {
	// initial validator
	validate = validator.New()

	// setup translator english language
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	var ok bool
	translator, ok = uni.GetTranslator("en")
	if !ok {
		panic("translator not found")
	}

	// register translator to validator
	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(err)
	}
}

// GetValidator return instance validator
func GetValidator() *validator.Validate {
	return validate
}

// GetTranslator return instance translator
func GetTranslator() ut.Translator {
	return translator
}
