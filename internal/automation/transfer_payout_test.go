package automation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"order/internal/payout"
	"order/internal/transfer"
)

func setupRouter() http.Handler {
	pricing := transfer.NewInMemoryRepository()
	pricing.SetRate("GBP", "NGN", 210)
	pricing.SetFee("GBP", "NGN", 1.99)

	transferSvc := transfer.NewService(pricing)
	transferCtrl := transfer.NewController(transferSvc)

	payoutRepo := payout.NewInMemoryRepository()
	payoutSvc := payout.NewService(payoutRepo, pricing)
	payoutCtrl := payout.NewController(payoutSvc)

	r := mux.NewRouter()
	transferCtrl.RegisterRoutes(r)
	payoutCtrl.RegisterRoutes(r)
	return r
}

func TestTransferCalculateValidationError(t *testing.T) {
	srv := httptest.NewServer(setupRouter())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/transfers/calculate", "application/json", bytes.NewBufferString(`{"source_currency":"GBP"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 got %d", resp.StatusCode)
	}
}

func TestPayoutInitiateUnsupportedCurrency(t *testing.T) {
	srv := httptest.NewServer(setupRouter())
	defer srv.Close()

	payload := payout.CreateRequest{
		SourceWalletID:     "wallet-1",
		DestinationBank:    "Bank",
		DestinationAccount: "1234567890",
		SourceCurrency:     "USD",
		TargetCurrency:     "NGN",
		SourceAmount:       10,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(srv.URL+"/payouts", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}
