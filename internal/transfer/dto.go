package transfer

// CalculateRequest encapsulates the payload for calculating a transfer quote.
type CalculateRequest struct {
	SourceCurrency string   `json:"source_currency" validate:"required,currency"`
	TargetCurrency string   `json:"target_currency" validate:"required,currency,nefield=SourceCurrency"`
	SourceAmount   float64  `json:"source_amount" validate:"required,gt=0"`
	FeeOverride    *float64 `json:"fee_override,omitempty" validate:"omitempty,gt=0,fee_lt_amount"`
}
