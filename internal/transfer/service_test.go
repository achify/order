package transfer

import (
	"context"
	"testing"
)

func TestCalculate(t *testing.T) {
	repo := NewInMemoryRepository()
	repo.SetRate("GBP", "NGN", 210)
	repo.SetFee("GBP", "NGN", 1.99)

	svc := NewService(repo)
	quote, err := svc.Calculate(context.Background(), CalculateRequest{SourceCurrency: "GBP", TargetCurrency: "NGN", SourceAmount: 50})
	if err != nil {
		t.Fatalf("calculate: %v", err)
	}
	if quote.TotalDebit != 51.99 {
		t.Fatalf("expected total debit 51.99, got %f", quote.TotalDebit)
	}
	if quote.TargetAmount != 10500 {
		t.Fatalf("expected target amount 10500, got %f", quote.TargetAmount)
	}
}
