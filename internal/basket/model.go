package basket

import "time"

// Basket holds items selected by a user before checkout
// TotalPrice is the sum of item price * quantity stored in cents

// Basket represents the overall basket entity

type Basket struct {
	ID         string    `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	AccountID  string    `db:"account_id" json:"account_id"`
	TotalPrice int64     `db:"total_price" json:"total_price"`
}

// Item represents an item within a basket with quantity and price per unit

type Item struct {
	BasketID     string `db:"basket_id" json:"basket_id"`
	ItemID       string `db:"item_id" json:"item_id"`
	Quantity     int    `db:"quantity" json:"quantity"`
	PricePerItem int64  `db:"price_per_item" json:"price_per_item"`
}
