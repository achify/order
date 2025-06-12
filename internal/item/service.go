package item

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

// Service handles item business logic.
type Service struct {
	Repo Repository
}

func NewService(r Repository) *Service { return &Service{Repo: r} }

func (s *Service) Create(ctx context.Context, dto ItemCreateDTO) (*Item, error) {
	it := &Item{
		ID:         ulid.Make().String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Name:       dto.Name,
		Price:      dto.Price,
		CategoryID: dto.CategoryID,
	}
	if err := s.Repo.Create(ctx, it); err != nil {
		return nil, err
	}
	return it, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Item, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Item, error) {
	return s.Repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id string, dto ItemUpdateDTO) (*Item, error) {
	it, err := s.Repo.GetByID(ctx, id)
	if err != nil || it == nil {
		return it, err
	}
	if dto.Name != nil {
		it.Name = *dto.Name
	}
	if dto.Price != nil {
		it.Price = *dto.Price
	}
	if dto.CategoryID != nil {
		it.CategoryID = *dto.CategoryID
	}
	it.UpdatedAt = time.Now()
	if err := s.Repo.Update(ctx, it); err != nil {
		return nil, err
	}
	return it, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}
