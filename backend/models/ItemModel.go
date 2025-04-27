package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model

	PoId          uint
	PurchaseOrder PurchaseOrder `gorm:"foreignKey:PoId"`
	ItemName      string        `gorm:"size:255" json:"item_name" validate:"required,min=3,max=255"`
	ItemCode      string        `gorm:"size:255;uniqueIndex" json:"item_code" validate:"required,min=3,max=255"`
	ItemPrice     float64       `gorm:"type:decimal(10,2)" json:"item_price" validate:"required"`
	ItemQty       int           `gorm:"size:255" json:"item_qty" validate:"required,min=1,max=255"`
}

func (i *Item) validate() error {
	validate := validator.New()
	return validate.Struct(i)
}
