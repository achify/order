package payout

import "time"

// Payout represents a payout transaction persisted by the service.
type Payout struct {
	ID                 string    `json:"id"`
	SourceWalletID     string    `json:"source_wallet_id"`
	DestinationBank    string    `json:"destination_bank"`
	DestinationAccount string    `json:"destination_account"`
	SourceCurrency     string    `json:"source_currency"`
	TargetCurrency     string    `json:"target_currency"`
	SourceAmount       float64   `json:"source_amount"`
	TargetAmount       float64   `json:"target_amount"`
	FeeAmount          float64   `json:"fee_amount"`
	ExchangeRate       float64   `json:"exchange_rate"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

const (
	// StatusPending indicates the payout has been created but not yet settled.
	StatusPending = "pending"
)
