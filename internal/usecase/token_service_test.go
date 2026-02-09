package usecase

import (
	"testing"
	"time"

	"github.com/example/validacion-pases/internal/security/auth"
)

func TestIssueTokenSuccess(t *testing.T) {
	users := auth.NewUserStore(map[string]string{"apiuser": "secret"})
	issuer, err := auth.NewTokenIssuer("HS256", "issuer", "aud", time.Hour, "hs-secret")
	if err != nil {
		t.Fatalf("issuer error: %v", err)
	}
	svc := NewTokenService(users, issuer)

	token, exp, err := svc.Issue("apiuser", "secret")
	if err != nil {
		t.Fatalf("issue error: %v", err)
	}
	if token == "" {
		t.Fatal("token empty")
	}
	if exp.IsZero() {
		t.Fatal("expiration not set")
	}
}

func TestIssueTokenInvalidCredentials(t *testing.T) {
	users := auth.NewUserStore(map[string]string{"apiuser": "secret"})
	issuer, _ := auth.NewTokenIssuer("HS256", "issuer", "aud", time.Hour, "hs-secret")
	svc := NewTokenService(users, issuer)

	_, _, err := svc.Issue("apiuser", "wrong")
	if err == nil {
		t.Fatal("expected error")
	}
}
