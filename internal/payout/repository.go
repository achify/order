package payout

import (
	"context"
	"sync"
)

// Repository abstracts persistence for payouts.
type Repository interface {
	Create(ctx context.Context, payout *Payout) error
}

// InMemoryRepository keeps payouts in memory, useful for tests.
type InMemoryRepository struct {
	mu      sync.RWMutex
	payouts map[string]*Payout
}

// NewInMemoryRepository creates an empty repository instance.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{payouts: make(map[string]*Payout)}
}

// Create persists the payout.
func (r *InMemoryRepository) Create(_ context.Context, payout *Payout) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.payouts[payout.ID] = payout
	return nil
}

// Get returns a payout by id for verification in tests.
func (r *InMemoryRepository) Get(id string) (*Payout, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.payouts[id]
	return p, ok
}
