package item

import "time"

// Item represents a product for sale
// Price is stored in cents to avoid floating point issues

type Item struct {
	ID         string    `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	Name       string    `db:"name" json:"name"`
	Price      int64     `db:"price" json:"price"`
	CategoryID string    `db:"category_id" json:"category_id"`
}
