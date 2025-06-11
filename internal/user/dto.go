package user

type UserCreateDTO struct {
	Username  string  `json:"username" validate:"required"`
	Password  string  `json:"password" validate:"required"`
	Roles     []Role  `json:"roles" validate:"required"`
	SellerID  *string `json:"seller_id,omitempty"`
	AccountID *string `json:"account_id,omitempty"`
}
