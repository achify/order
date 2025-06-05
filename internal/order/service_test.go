package order

import (
	"context"
	"testing"
)

type mockRepo struct {
	store map[string]*Order
}

func newMock() *mockRepo { return &mockRepo{store: make(map[string]*Order)} }

func (m *mockRepo) Create(ctx context.Context, o *Order) error {
	m.store[o.ID] = o
	return nil
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (*Order, error) {
	if o, ok := m.store[id]; ok {
		return o, nil
	}
	return nil, nil
}
func (m *mockRepo) List(ctx context.Context, deliveryID string) ([]Order, error) { return nil, nil }
func (m *mockRepo) Update(ctx context.Context, o *Order) error                   { m.store[o.ID] = o; return nil }
func (m *mockRepo) Delete(ctx context.Context, id string) error                  { delete(m.store, id); return nil }

func TestServiceCreateAndGet(t *testing.T) {
	repo := newMock()
	svc := NewService(repo)
	o, err := svc.Create(context.Background(), OrderCreateDTO{
		ReceiverID: "r", AccountID: "a", SellerID: "s", DeliveryID: "d", BasketID: "b",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := svc.Get(context.Background(), o.ID)
	if err != nil || got == nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != o.ID {
		t.Fatalf("want %s got %s", o.ID, got.ID)
	}
}
