package auth

import "testing"

func TestGenerateAndRefresh(t *testing.T) {
	svc := NewService([]byte("test"))
	tok, refresh, err := svc.GenerateToken("user", []string{"admin"})
	if err != nil || tok == "" || refresh == "" {
		t.Fatalf("generate: %v", err)
	}
	newTok, newRefresh, err := svc.Refresh(refresh)
	if err != nil || newTok == "" || newRefresh == "" {
		t.Fatalf("refresh: %v", err)
	}
}
