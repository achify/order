package router

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"order/internal/auth"
	"order/internal/basket"
	"order/internal/item"
	ord "order/internal/order"
	"order/internal/payout"
	"order/internal/transfer"
	usr "order/internal/user"
)

// mockRepo is a simple in-memory Repository implementation for tests
type mockRepo struct{ store map[string]*ord.Order }
type itemRepo struct{}
type basketRepo struct{}

func newMockRepo() *mockRepo { return &mockRepo{store: make(map[string]*ord.Order)} }

func (m *mockRepo) Create(_ context.Context, o *ord.Order) error { m.store[o.ID] = o; return nil }
func (m *mockRepo) GetByID(_ context.Context, id string) (*ord.Order, error) {
	if o, ok := m.store[id]; ok {
		return o, nil
	}
	return nil, nil
}
func (m *mockRepo) List(context.Context, string) ([]ord.Order, error) {
	var list []ord.Order
	return list, nil
}
func (m *mockRepo) Update(context.Context, *ord.Order) error            { return nil }
func (m *mockRepo) Delete(_ context.Context, id string) error           { delete(m.store, id); return nil }
func (i *itemRepo) Create(context.Context, *item.Item) error            { return nil }
func (i *itemRepo) GetByID(context.Context, string) (*item.Item, error) { return nil, nil }
func (i *itemRepo) List(context.Context) ([]item.Item, error)           { return nil, nil }
func (i *itemRepo) Update(context.Context, *item.Item) error            { return nil }
func (i *itemRepo) Delete(context.Context, string) error                { return nil }

func (b *basketRepo) Create(context.Context, *basket.Basket) error             { return nil }
func (b *basketRepo) GetByID(context.Context, string) (*basket.Basket, error)  { return nil, nil }
func (b *basketRepo) List(context.Context) ([]basket.Basket, error)            { return nil, nil }
func (b *basketRepo) Update(context.Context, *basket.Basket) error             { return nil }
func (b *basketRepo) Delete(context.Context, string) error                     { return nil }
func (b *basketRepo) AddItem(context.Context, *basket.Item) error              { return nil }
func (b *basketRepo) UpdateItem(context.Context, *basket.Item) error           { return nil }
func (b *basketRepo) DeleteItem(context.Context, string, string) error         { return nil }
func (b *basketRepo) ListItems(context.Context, string) ([]basket.Item, error) { return nil, nil }

// memUserRepo is a simple in-memory user repository.
type memUserRepo struct{ user *usr.User }

func (m *memUserRepo) Create(_ context.Context, u *usr.User) error { m.user = u; return nil }
func (m *memUserRepo) GetByUsername(_ context.Context, _ string) (*usr.User, error) {
	return m.user, nil
}
func (m *memUserRepo) Update(_ context.Context, _ *usr.User) error { return nil }

func TestCreateOrderStatusCode(t *testing.T) {
	repo := newMockRepo()
	svc := ord.NewService(repo)
	ctrl := ord.NewController(svc)
	itemCtrl := item.NewController(item.NewService(&itemRepo{}))
	basketCtrl := basket.NewController(basket.NewService(&basketRepo{}))

	secret := []byte("secret")
	authSvc := auth.NewService(secret)
	userRepo := &memUserRepo{}
	userSvc := usr.NewService(userRepo)
	// seed admin user
	_, _ = userSvc.Create(context.Background(), usr.UserCreateDTO{Username: "admin", Password: "pass", Roles: []usr.Role{usr.RoleAdmin}})
	authCtrl := auth.NewController(authSvc, userSvc)
	userCtrl := usr.NewController(userSvc)

	pricing := transfer.NewInMemoryRepository()
	pricing.SetRate("GBP", "NGN", 210)
	transferSvc := transfer.NewService(pricing)
	transferCtrl := transfer.NewController(transferSvc)

	payoutRepo := payout.NewInMemoryRepository()
	payoutSvc := payout.NewService(payoutRepo, pricing)
	payoutCtrl := payout.NewController(payoutSvc)

	r := New(ctrl, itemCtrl, basketCtrl, secret, authCtrl, userCtrl, transferCtrl, payoutCtrl)
	ts := httptest.NewServer(r)
	defer ts.Close()

	token, _, err := authSvc.GenerateToken("admin", []string{string(usr.RoleAdmin)})
	if err != nil {
		t.Fatalf("token: %v", err)
	}

	body := []byte(`{"receiver_id":"r","account_id":"a","seller_id":"s","delivery_id":"d","basket_id":"b"}`)
	req, _ := http.NewRequest("POST", ts.URL+"/orders", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want %d got %d", http.StatusCreated, resp.StatusCode)
	}
}
