package order

// OrderCreateDTO defines input for creating an order.
// Validation tags ensure data integrity

type OrderCreateDTO struct {
	ReceiverID string `json:"receiver_id" validate:"required"`
	AccountID  string `json:"account_id" validate:"required"`
	SellerID   string `json:"seller_id" validate:"required"`
	DeliveryID string `json:"delivery_id" validate:"required"`
	BasketID   string `json:"basket_id" validate:"required"`
	Status     string `json:"status"`
}

// OrderUpdateDTO defines input for updating an order
// All fields are optional for PATCH

type OrderUpdateDTO struct {
	ReceiverID *string `json:"receiver_id"`
	AccountID  *string `json:"account_id"`
	SellerID   *string `json:"seller_id"`
	DeliveryID *string `json:"delivery_id"`
	BasketID   *string `json:"basket_id"`
	Status     *string `json:"status"`
}
