package order

import "time"

// Order domain model
// contains basic fields for an e-commerce order

type Order struct {
	ID         string    `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	ReceiverID string    `db:"receiver_id" json:"receiver_id"`
	AccountID  string    `db:"account_id" json:"account_id"`
	SellerID   string    `db:"seller_id" json:"seller_id"`
	DeliveryID string    `db:"delivery_id" json:"delivery_id"`
	BasketID   string    `db:"basket_id" json:"basket_id"`
}
