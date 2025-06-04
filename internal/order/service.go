package order

import (
    "context"
    "time"

    "github.com/oklog/ulid/v2"
)

// Service handles business logic for orders

type Service struct {
    Repo Repository
}

func NewService(r Repository) *Service { return &Service{Repo: r} }

func (s *Service) Create(ctx context.Context, dto OrderCreateDTO) (*Order, error) {
    o := &Order{
        ID:         ulid.Make().String(),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        ReceiverID: dto.ReceiverID,
        AccountID:  dto.AccountID,
        SellerID:   dto.SellerID,
        DeliveryID: dto.DeliveryID,
        BasketID:   dto.BasketID,
    }
    if err := s.Repo.Create(ctx, o); err != nil {
        return nil, err
    }
    return o, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Order, error) {
    return s.Repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Order, error) {
    return s.Repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id string, dto OrderUpdateDTO) (*Order, error) {
    o, err := s.Repo.GetByID(ctx, id)
    if err != nil || o == nil {
        return o, err
    }
    if dto.ReceiverID != nil {
        o.ReceiverID = *dto.ReceiverID
    }
    if dto.AccountID != nil {
        o.AccountID = *dto.AccountID
    }
    if dto.SellerID != nil {
        o.SellerID = *dto.SellerID
    }
    if dto.DeliveryID != nil {
        o.DeliveryID = *dto.DeliveryID
    }
    if dto.BasketID != nil {
        o.BasketID = *dto.BasketID
    }
    o.UpdatedAt = time.Now()
    if err := s.Repo.Update(ctx, o); err != nil {
        return nil, err
    }
    return o, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
    return s.Repo.Delete(ctx, id)
}
