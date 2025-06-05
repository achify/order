package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
)

// Service encapsulates business logic for deliveries.
type Service struct {
	Repo       Repository
	HTTPClient *http.Client
}

func NewService(r Repository) *Service {
	return &Service{Repo: r, HTTPClient: &http.Client{Timeout: 10 * time.Second}}
}

func (s *Service) Create(ctx context.Context, dto CreateDTO) (*Delivery, error) {
	d := &Delivery{
		ID:           ulid.Make().String(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Provider:     dto.Provider,
		TrackingCode: dto.TrackingCode,
		Status:       dto.Status,
	}
	if err := s.Repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Delivery, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Delivery, error) {
	return s.Repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id string, dto UpdateDTO) (*Delivery, error) {
	d, err := s.Repo.GetByID(ctx, id)
	if err != nil || d == nil {
		return d, err
	}
	if dto.Status != nil {
		d.Status = *dto.Status
	}
	d.UpdatedAt = time.Now()
	if err := s.Repo.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

// SyncLocations fetches pickup locations from Omniva and stores them.
func (s *Service) SyncLocations(ctx context.Context) error {
	resp, err := s.HTTPClient.Get("https://www.omniva.ee/locations.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	now := time.Now()
	for _, item := range data {
		b, _ := json.Marshal(item)
		loc := &Location{ID: item["ZIP"].(string), Provider: "omniva", Data: string(b), UpdatedAt: now}
		if err := s.Repo.UpsertLocation(ctx, loc); err != nil {
			return err
		}
	}
	return nil
}
