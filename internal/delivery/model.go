package delivery

import "time"

// Delivery represents a shipment tracked by the system.
type Delivery struct {
	ID           string    `db:"id" json:"id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Provider     string    `db:"provider" json:"provider"`
	TrackingCode string    `db:"tracking_code" json:"tracking_code"`
	Status       string    `db:"status" json:"status"`
}

// Location represents a delivery pick-up location (e.g. parcel machine).
type Location struct {
	ID        string    `db:"id" json:"id"`
	Provider  string    `db:"provider" json:"provider"`
	Data      string    `db:"data" json:"data"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
