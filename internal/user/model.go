package user

import "time"

// Role defines access level of a user
//
// Additional roles may be added to support fine grained RBAC.
type Role string

const (
	RoleAdmin            Role = "admin"
	RoleSupportManager   Role = "support_manager"
	RoleSeller           Role = "seller"
	RoleCustomer         Role = "customer"
	RoleLogisticsManager Role = "logistics_manager"
	RolePaymentManager   Role = "payment_manager"
	RoleAuditor          Role = "auditor"
)

// User represents application user with login credentials and profile link.
type User struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"`
	Roles     []Role    `db:"roles" json:"roles"`
	Active    bool      `db:"active" json:"active"`
	SellerID  *string   `db:"seller_id" json:"seller_id,omitempty"`
	AccountID *string   `db:"account_id" json:"account_id,omitempty"`
}
