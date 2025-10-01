package transfer

import "time"

// Quote describes the calculated transfer details before execution.
type Quote struct {
	ID             string    `json:"id"`
	SourceCurrency string    `json:"source_currency"`
	TargetCurrency string    `json:"target_currency"`
	SourceAmount   float64   `json:"source_amount"`
	TargetAmount   float64   `json:"target_amount"`
	ExchangeRate   float64   `json:"exchange_rate"`
	FeeAmount      float64   `json:"fee_amount"`
	TotalDebit     float64   `json:"total_debit"`
	CreatedAt      time.Time `json:"created_at"`
}
