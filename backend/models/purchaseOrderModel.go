package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PurchaseOrder struct {
	gorm.Model

	CustomerName    string  `gorm:"size:255" json:"customer_name" validate:"required,min=3,max=255"`
	CustomerAddress string  `gorm:"size:255" json:"shipping_address" validate:"required,min=10,max=255"`
	PoNumber        string  `gorm:"size:255;uniqueIndex" json:"po_number" validate:"required,min=3,max=255"`
	TotalAmount     float64 `gorm:"type:decimal(10,2)" json:"total_amount" validate:"required"`
	ShippingDate    string  `gorm:"size:255" json:"shipping_date" validate:"required,min=3,max=255"`
}

func (po *PurchaseOrder) validate() error {
	validate := validator.New()
	return validate.Struct(po)
}
