package payment

import (
	"context"
	ord "order/internal/order"
	"testing"
)

type mockOrderRepo struct{ store map[string]*ord.Order }

func newOrderMock() *mockOrderRepo { return &mockOrderRepo{store: make(map[string]*ord.Order)} }

func (m *mockOrderRepo) Create(ctx context.Context, o *ord.Order) error {
	m.store[o.ID] = o
	return nil
}
func (m *mockOrderRepo) GetByID(ctx context.Context, id string) (*ord.Order, error) {
	if o, ok := m.store[id]; ok {
		return o, nil
	}
	return nil, nil
}
func (m *mockOrderRepo) List(ctx context.Context, deliveryID string) ([]ord.Order, error) {
	return nil, nil
}
func (m *mockOrderRepo) Update(ctx context.Context, o *ord.Order) error {
	m.store[o.ID] = o
	return nil
}
func (m *mockOrderRepo) Delete(ctx context.Context, id string) error { delete(m.store, id); return nil }

type mockRepo struct{ store map[string]*Payment }

func newMock() *mockRepo { return &mockRepo{store: make(map[string]*Payment)} }

func (m *mockRepo) Create(ctx context.Context, p *Payment) error { m.store[p.ID] = p; return nil }
func (m *mockRepo) GetByID(ctx context.Context, id string) (*Payment, error) {
	if p, ok := m.store[id]; ok {
		return p, nil
	}
	return nil, nil
}

func TestAcceptPayment(t *testing.T) {
	pRepo := newMock()
	oRepo := newOrderMock()
	order := &ord.Order{ID: "1", ReceiverID: "r", AccountID: "a", SellerID: "s", DeliveryID: "d", BasketID: "b", Status: ord.StatusNew}
	oRepo.Create(context.Background(), order)

	svc := NewService(pRepo, oRepo)
	pay, err := svc.Accept(context.Background(), CreateDTO{OrderID: "1", Amount: 100})
	if err != nil {
		t.Fatalf("accept: %v", err)
	}
	if pay.Status != "paid" {
		t.Fatalf("want paid got %s", pay.Status)
	}
	if oRepo.store["1"].Status != ord.StatusPaid {
		t.Fatalf("order not updated")
	}
}
