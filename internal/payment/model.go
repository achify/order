package payment

import "time"

// Payment represents a received payment for an order.
type Payment struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	OrderID   string    `db:"order_id" json:"order_id"`
	Amount    int64     `db:"amount" json:"amount"`
	Status    string    `db:"status" json:"status"`
}
