package payout

import (
	"context"
	"testing"

	"order/internal/transfer"
)

func TestInitiate(t *testing.T) {
	pricing := transfer.NewInMemoryRepository()
	pricing.SetRate("GBP", "NGN", 210)
	pricing.SetFee("GBP", "NGN", 1.99)

	repo := NewInMemoryRepository()
	svc := NewService(repo, pricing)

	req := CreateRequest{
		SourceWalletID:     "wallet-1",
		DestinationBank:    "Nigerian Bank",
		DestinationAccount: "1234567890",
		SourceCurrency:     "GBP",
		TargetCurrency:     "NGN",
		SourceAmount:       50,
	}

	payout, err := svc.Initiate(context.Background(), req)
	if err != nil {
		t.Fatalf("initiate: %v", err)
	}
	if payout.Status != StatusPending {
		t.Fatalf("expected pending status, got %s", payout.Status)
	}
	if payout.TargetAmount != 10500 {
		t.Fatalf("expected NGN 10500, got %f", payout.TargetAmount)
	}
}
