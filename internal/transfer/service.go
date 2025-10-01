package transfer

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
)

// Service performs business operations for transfer calculations.
type Service struct {
	repo Repository
}

// NewService builds a Service instance.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Calculate produces a transfer quote using configured pricing metadata.
func (s *Service) Calculate(ctx context.Context, req CalculateRequest) (*Quote, error) {
	var (
		rate float64
		fee  float64
	)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		rate, err = s.repo.GetExchangeRate(ctx, req.SourceCurrency, req.TargetCurrency)
		return err
	})

	g.Go(func() error {
		if req.FeeOverride != nil {
			fee = *req.FeeOverride
			return nil
		}
		f, err := s.repo.GetFlatFee(ctx, req.SourceCurrency, req.TargetCurrency)
		if err != nil {
			if err == ErrFeeNotFound {
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
		// Default micro-transaction fee when none configured explicitly.
		fee = 1.99
	}

	quote := &Quote{
		ID:             ulid.Make().String(),
		SourceCurrency: req.SourceCurrency,
		TargetCurrency: req.TargetCurrency,
		SourceAmount:   req.SourceAmount,
		ExchangeRate:   rate,
		FeeAmount:      fee,
		TotalDebit:     req.SourceAmount + fee,
		TargetAmount:   rate * req.SourceAmount,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.repo.SaveQuote(ctx, quote); err != nil {
		return nil, err
	}
	return quote, nil
}
