package basket

import (
	"context"
	"testing"
)

type mockRepo struct{ store map[string]*Basket }

func newMock() *mockRepo { return &mockRepo{store: make(map[string]*Basket)} }

func (m *mockRepo) Create(_ context.Context, b *Basket) error { m.store[b.ID] = b; return nil }
func (m *mockRepo) GetByID(_ context.Context, id string) (*Basket, error) {
	if b, ok := m.store[id]; ok {
		return b, nil
	}
	return nil, nil
}
func (m *mockRepo) List(context.Context) ([]Basket, error)            { return nil, nil }
func (m *mockRepo) Update(_ context.Context, b *Basket) error         { m.store[b.ID] = b; return nil }
func (m *mockRepo) Delete(_ context.Context, id string) error         { delete(m.store, id); return nil }
func (m *mockRepo) AddItem(context.Context, *Item) error              { return nil }
func (m *mockRepo) UpdateItem(context.Context, *Item) error           { return nil }
func (m *mockRepo) DeleteItem(context.Context, string, string) error  { return nil }
func (m *mockRepo) ListItems(context.Context, string) ([]Item, error) { return nil, nil }

func TestServiceCreateAndGet(t *testing.T) {
	repo := newMock()
	svc := NewService(repo)
	b, err := svc.Create(context.Background(), BasketCreateDTO{AccountID: "a"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := svc.Get(context.Background(), b.ID)
	if err != nil || got == nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != b.ID {
		t.Fatalf("want %s got %s", b.ID, got.ID)
	}
}
