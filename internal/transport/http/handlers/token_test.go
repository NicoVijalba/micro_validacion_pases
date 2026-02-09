package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/validacion-pases/internal/security/auth"
	"github.com/example/validacion-pases/internal/usecase"
)

func TestIssueTokenHandler(t *testing.T) {
	users := auth.NewUserStore(map[string]string{"apiuser": "secret"})
	issuer, _ := auth.NewTokenIssuer("HS256", "issuer", "aud", time.Hour, "hs-secret")
	h := NewTokenHandler(usecase.NewTokenService(users, issuer))

	r := httptest.NewRequest(http.MethodPost, "/v1/token", bytes.NewReader([]byte(`{"username":"apiuser","password":"secret"}`)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Issue(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
