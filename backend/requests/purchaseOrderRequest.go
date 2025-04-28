package requests

type PurchaseOrderItemRequest struct {
	ItemName  string  `json:"item_name" validate:"required"`
	ItemCode  string  `json:"item_code" validate:"required"`
	ItemPrice float64 `json:"item_price" validate:"required"`
	ItemQty   int     `json:"item_qty" validate:"required"`
}

type PurchaseOrderRequest struct {
	CustomerName    string                     `json:"customer_name" validate:"required"`
	CustomerAddress string                     `json:"customer_address" validate:"required"`
	PoNumber        string                     `json:"po_number" validate:"required"`
	TotalAmount     float64                    `json:"total_amount" validate:"required"`
	ShippingDate    string                     `json:"shipping_date" validate:"required"`
	Items           []PurchaseOrderItemRequest `json:"items" validate:"required,dive"`
}
