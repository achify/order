package payout

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"

	"order/internal/transfer"
)

// Service coordinates payout transactions.
type Service struct {
	repo    Repository
	pricing transfer.PricingProvider
}

// NewService constructs a payout Service instance.
func NewService(repo Repository, pricing transfer.PricingProvider) *Service {
	return &Service{repo: repo, pricing: pricing}
}

// Initiate creates a payout transaction ready for downstream processing.
func (s *Service) Initiate(ctx context.Context, req CreateRequest) (*Payout, error) {
	var (
		rate float64
		fee  float64
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		rate, err = s.pricing.GetExchangeRate(ctx, req.SourceCurrency, req.TargetCurrency)
		return err
	})

	g.Go(func() error {
		f, err := s.pricing.GetFlatFee(ctx, req.SourceCurrency, req.TargetCurrency)
		if err != nil {
			if err == transfer.ErrFeeNotFound {
				fee = 0
				return nil
			}
			return err
		}
		fee = f
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if req.SourceAmount <= 50 && fee == 0 {
		fee = 1.99
	}

	payout := &Payout{
		ID:                 ulid.Make().String(),
		SourceWalletID:     req.SourceWalletID,
		DestinationBank:    req.DestinationBank,
		DestinationAccount: req.DestinationAccount,
		SourceCurrency:     req.SourceCurrency,
		TargetCurrency:     req.TargetCurrency,
		SourceAmount:       req.SourceAmount,
		ExchangeRate:       rate,
		FeeAmount:          fee,
		TargetAmount:       rate * req.SourceAmount,
		Status:             StatusPending,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, payout); err != nil {
		return nil, err
	}

	return payout, nil
}
