package basket

// BasketCreateDTO defines input for creating a basket.
type BasketCreateDTO struct {
	AccountID string `json:"account_id" validate:"required"`
}

// BasketUpdateDTO defines input for updating a basket.
type BasketUpdateDTO struct {
	AccountID  *string `json:"account_id"`
	TotalPrice *int64  `json:"total_price"`
}

// ItemDTO is used when manipulating basket items.
type ItemDTO struct {
	ItemID       string `json:"item_id" validate:"required"`
	Quantity     int    `json:"quantity" validate:"required"`
	PricePerItem int64  `json:"price_per_item" validate:"required"`
}
