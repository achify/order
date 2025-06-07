package payment

// CreateDTO describes a payment acceptance request.
type CreateDTO struct {
	OrderID string `json:"order_id" validate:"required"`
	Amount  int64  `json:"amount" validate:"required"`
}
