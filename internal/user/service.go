package user

import (
	"context"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

// Service encapsulates business logic around users.
type Service struct {
	Repo Repository
}

func NewService(r Repository) *Service { return &Service{Repo: r} }

func (s *Service) Create(ctx context.Context, dto UserCreateDTO) (*User, error) {
	u := &User{
		ID:        ulid.Make().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  dto.Username,
		Password:  dto.Password,
		Roles:     dto.Roles,
		Active:    true,
		SellerID:  dto.SellerID,
		AccountID: dto.AccountID,
	}
	if err := s.Repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// Authenticate verifies credentials and returns the user if valid.
func (s *Service) Authenticate(ctx context.Context, username, password string) (*User, error) {
	u, err := s.Repo.GetByUsername(ctx, username)
	if err != nil || u == nil || !u.Active || u.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}
