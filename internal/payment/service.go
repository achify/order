package payment

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	ord "order/internal/order"
)

// Service implements payment acceptance logic.
type Service struct {
	Repo      Repository
	OrderRepo ord.Repository
}

func NewService(r Repository, or ord.Repository) *Service {
	return &Service{Repo: r, OrderRepo: or}
}

func (s *Service) Accept(ctx context.Context, dto CreateDTO) (*Payment, error) {
	p := &Payment{
		ID:        ulid.Make().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		OrderID:   dto.OrderID,
		Amount:    dto.Amount,
		Status:    "paid",
	}
	if err := s.Repo.Create(ctx, p); err != nil {
		return nil, err
	}
	if o, err := s.OrderRepo.GetByID(ctx, dto.OrderID); err == nil && o != nil {
		o.Status = ord.StatusPaid
		o.UpdatedAt = time.Now()
		if err := s.OrderRepo.Update(ctx, o); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Payment, error) {
	return s.Repo.GetByID(ctx, id)
}
