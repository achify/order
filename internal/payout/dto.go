package payout

// CreateRequest models the payload required to initiate a payout.
type CreateRequest struct {
	SourceWalletID     string  `json:"source_wallet_id" validate:"required"`
	DestinationBank    string  `json:"destination_bank" validate:"required,min=2"`
	DestinationAccount string  `json:"destination_account" validate:"required,account"`
	SourceCurrency     string  `json:"source_currency" validate:"required,currency"`
	TargetCurrency     string  `json:"target_currency" validate:"required,currency,nefield=SourceCurrency"`
	SourceAmount       float64 `json:"source_amount" validate:"required,gt=0"`
}
