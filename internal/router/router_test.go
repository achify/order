package router

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"order/internal/auth"
	ord "order/internal/order"
	usr "order/internal/user"
)

// mockRepo is a simple in-memory Repository implementation for tests
type mockRepo struct{ store map[string]*ord.Order }

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
func (m *mockRepo) Update(context.Context, *ord.Order) error  { return nil }
func (m *mockRepo) Delete(_ context.Context, id string) error { delete(m.store, id); return nil }

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

	secret := []byte("secret")
	authSvc := auth.NewService(secret)
	userRepo := &memUserRepo{}
	userSvc := usr.NewService(userRepo)
	// seed admin user
	_, _ = userSvc.Create(context.Background(), usr.UserCreateDTO{Username: "admin", Password: "pass", Roles: []usr.Role{usr.RoleAdmin}})
	authCtrl := auth.NewController(authSvc, userSvc)
	userCtrl := usr.NewController(userSvc)

	r := New(ctrl, secret, authCtrl, userCtrl)
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
