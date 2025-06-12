package basket

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

// Service handles basket business logic.
type Service struct{ Repo Repository }

func NewService(r Repository) *Service { return &Service{Repo: r} }

func (s *Service) Create(ctx context.Context, dto BasketCreateDTO) (*Basket, error) {
	b := &Basket{
		ID:         ulid.Make().String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		AccountID:  dto.AccountID,
		TotalPrice: 0,
	}
	if err := s.Repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Basket, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Basket, error) {
	return s.Repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id string, dto BasketUpdateDTO) (*Basket, error) {
	b, err := s.Repo.GetByID(ctx, id)
	if err != nil || b == nil {
		return b, err
	}
	if dto.AccountID != nil {
		b.AccountID = *dto.AccountID
	}
	if dto.TotalPrice != nil {
		b.TotalPrice = *dto.TotalPrice
	}
	b.UpdatedAt = time.Now()
	if err := s.Repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *Service) AddItem(ctx context.Context, basketID string, dto ItemDTO) error {
	it := &Item{BasketID: basketID, ItemID: dto.ItemID, Quantity: dto.Quantity, PricePerItem: dto.PricePerItem}
	return s.Repo.AddItem(ctx, it)
}

func (s *Service) UpdateItem(ctx context.Context, basketID, itemID string, dto ItemDTO) error {
	it := &Item{BasketID: basketID, ItemID: itemID, Quantity: dto.Quantity, PricePerItem: dto.PricePerItem}
	return s.Repo.UpdateItem(ctx, it)
}

func (s *Service) DeleteItem(ctx context.Context, basketID, itemID string) error {
	return s.Repo.DeleteItem(ctx, basketID, itemID)
}

func (s *Service) ListItems(ctx context.Context, basketID string) ([]Item, error) {
	return s.Repo.ListItems(ctx, basketID)
}
