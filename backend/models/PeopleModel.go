package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type People struct {
	gorm.Model

	Name	string `gorm:"size:255" json:"name" validate:"required,min=3,max=255"`
	Address string `gorm:"size:255" json:"address" validate:"required,min=3,max=255"`
	Phone	string `gorm:"size:255" json:"phone" validate:"required,min=3,max=255"`
}

func (p *People) validate() error {
	validate := validator.New()
	return validate.Struct(p)
}