package transfer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// PricingProvider exposes read operations used for pricing calculations.
type PricingProvider interface {
	GetExchangeRate(ctx context.Context, sourceCurrency, targetCurrency string) (float64, error)
	GetFlatFee(ctx context.Context, sourceCurrency, targetCurrency string) (float64, error)
}

// Repository represents the storage for transfer quotes and pricing metadata.
type Repository interface {
	PricingProvider
	SaveQuote(ctx context.Context, quote *Quote) error
}

var (
	// ErrRateNotFound is returned when an exchange rate is missing for the provided currency pair.
	ErrRateNotFound = errors.New("exchange rate not found")
	// ErrFeeNotFound is returned when a flat fee configuration is missing.
	ErrFeeNotFound = errors.New("fee not found")
)

// InMemoryRepository is a thread-safe repository backed by maps for tests and local development.
type InMemoryRepository struct {
	mu     sync.RWMutex
	rates  map[string]float64
	fees   map[string]float64
	quotes map[string]*Quote
}

// NewInMemoryRepository constructs an empty in-memory repository instance.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		rates:  make(map[string]float64),
		fees:   make(map[string]float64),
		quotes: make(map[string]*Quote),
	}
}

// SetRate seeds an exchange rate for the provided currency pair.
func (r *InMemoryRepository) SetRate(sourceCurrency, targetCurrency string, rate float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rates[pairKey(sourceCurrency, targetCurrency)] = rate
}

// SetFee seeds a flat fee for the provided currency pair.
func (r *InMemoryRepository) SetFee(sourceCurrency, targetCurrency string, fee float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fees[pairKey(sourceCurrency, targetCurrency)] = fee
}

// GetExchangeRate returns the configured rate for the currency pair.
func (r *InMemoryRepository) GetExchangeRate(_ context.Context, sourceCurrency, targetCurrency string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := pairKey(sourceCurrency, targetCurrency)
	rate, ok := r.rates[key]
	if !ok {
		return 0, ErrRateNotFound
	}
	return rate, nil
}

// GetFlatFee returns the configured flat fee for the currency pair.
func (r *InMemoryRepository) GetFlatFee(_ context.Context, sourceCurrency, targetCurrency string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := pairKey(sourceCurrency, targetCurrency)
	fee, ok := r.fees[key]
	if !ok {
		return 0, ErrFeeNotFound
	}
	return fee, nil
}

// SaveQuote stores a generated quote for retrieval.
func (r *InMemoryRepository) SaveQuote(_ context.Context, quote *Quote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.quotes[quote.ID] = quote
	return nil
}

// GetQuote retrieves a previously generated quote by ID.
func (r *InMemoryRepository) GetQuote(id string) (*Quote, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q, ok := r.quotes[id]
	return q, ok
}

func pairKey(sourceCurrency, targetCurrency string) string {
	return fmt.Sprintf("%s:%s", sourceCurrency, targetCurrency)
}
