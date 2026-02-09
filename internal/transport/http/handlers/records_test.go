package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/validacion-pases/internal/domain"
	"github.com/example/validacion-pases/internal/security/auth"
	"github.com/example/validacion-pases/internal/transport/http/middleware"
	"github.com/example/validacion-pases/internal/usecase"
)

type testRepo struct{}

func (testRepo) Insert(_ context.Context, _ domain.Record) (int64, error) { return 123, nil }

func TestCreateRecordHandler(t *testing.T) {
	h := NewRecordHandler(usecase.NewRecordService(testRepo{}))
	body := []byte(`{"nave":"NAVE 1","viaje":"VJ1","cliente":"CLIENTE 1","booking":"BK1","rama":"internacional","contenedor_serie":"ABCU1234567","fecha_real":"2026-02-09","dias_libre":2,"puerto_descargue":"Balboa"}`)
	r := httptest.NewRequest(http.MethodPost, "/v1/records", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	claims := &auth.Claims{Subject: "user-1"}
	r = r.WithContext(middleware.WithClaims(r.Context(), claims))

	w := httptest.NewRecorder()
	h.Create(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateRecordBadPayload(t *testing.T) {
	h := NewRecordHandler(usecase.NewRecordService(testRepo{}))
	r := httptest.NewRequest(http.MethodPost, "/v1/records", bytes.NewReader([]byte(`{"bad":1}`)))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(middleware.WithClaims(r.Context(), &auth.Claims{Subject: "user-1"}))

	w := httptest.NewRecorder()
	h.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
