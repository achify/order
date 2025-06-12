package order

import "time"

// Order domain model
// contains basic fields for an e-commerce order

type Order struct {
	ID         string      `db:"id" json:"id"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
	ReceiverID string      `db:"receiver_id" json:"receiver_id"`
	AccountID  string      `db:"account_id" json:"account_id"`
	SellerID   string      `db:"seller_id" json:"seller_id"`
	DeliveryID string      `db:"delivery_id" json:"delivery_id"`
	BasketID   string      `db:"basket_id" json:"basket_id"`
	Status     OrderStatus `db:"status" json:"status"`
}

// OrderStatus represents current lifecycle state of an order.
type OrderStatus string

const (
	StatusNew              OrderStatus = "new"
	StatusPendingPayment   OrderStatus = "pending_payment"
	StatusPaid             OrderStatus = "paid"
	StatusProcessing       OrderStatus = "processing"
	StatusInAssembly       OrderStatus = "in_assembly"
	StatusReadyForShipment OrderStatus = "ready_for_shipment"
	StatusHandedOver       OrderStatus = "handed_over_for_delivery"
	StatusInTransit        OrderStatus = "in_transit"
	StatusDelivered        OrderStatus = "delivered"
	StatusCompleted        OrderStatus = "completed"
	StatusReturnRequested  OrderStatus = "return_requested"
	StatusReturned         OrderStatus = "returned"
	StatusCancelled        OrderStatus = "cancelled"
	StatusRefunded         OrderStatus = "refunded"
)
