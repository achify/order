package item

import (
	"context"
	"testing"
)

type mockRepo struct{ store map[string]*Item }

func newMock() *mockRepo { return &mockRepo{store: make(map[string]*Item)} }

func (m *mockRepo) Create(_ context.Context, it *Item) error { m.store[it.ID] = it; return nil }
func (m *mockRepo) GetByID(_ context.Context, id string) (*Item, error) {
	if it, ok := m.store[id]; ok {
		return it, nil
	}
	return nil, nil
}
func (m *mockRepo) List(context.Context) ([]Item, error)      { return nil, nil }
func (m *mockRepo) Update(_ context.Context, it *Item) error  { m.store[it.ID] = it; return nil }
func (m *mockRepo) Delete(_ context.Context, id string) error { delete(m.store, id); return nil }

func TestServiceCreateAndGet(t *testing.T) {
	repo := newMock()
	svc := NewService(repo)
	it, err := svc.Create(context.Background(), ItemCreateDTO{Name: "ipad", Price: 100, CategoryID: "c"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := svc.Get(context.Background(), it.ID)
	if err != nil || got == nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != it.ID {
		t.Fatalf("want %s got %s", it.ID, got.ID)
	}
}
