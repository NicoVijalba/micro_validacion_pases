package handlers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/validacion-pases/internal/domain"
	"github.com/example/validacion-pases/internal/security/auth"
	"github.com/example/validacion-pases/internal/transport/http/middleware"
	"github.com/example/validacion-pases/internal/usecase"
)

type testRepo struct{}

func (testRepo) Insert(_ context.Context, _ domain.Record) (int64, error) { return 123, nil }
func (testRepo) FindByID(_ context.Context, id int64) (domain.Record, error) {
	return domain.Record{
		ID:                  id,
		Emision:             time.Date(2026, 2, 17, 9, 41, 45, 0, time.UTC),
		Nave:                "NYK DENEB",
		Viaje:               "072E",
		Cliente:             "CAPITAL PACIFICO, S.A.",
		Booking:             "YMLUL160382911",
		Rama:                "internacional",
		Contenedor:          "YMLU5374938",
		PuertoDescargue:     "RODMAN",
		LibreRetencionHasta: time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC),
		DiasLibre:           17,
		Transportista:       "",
		TituloTerminal:      "PANAMA PORTS COMPANY (RODMAN)",
		UsuarioFirma:        "Admin",
		CreatedAt:           time.Date(2026, 2, 17, 9, 41, 45, 0, time.UTC),
	}, nil
}

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

func TestCreateRecordAcceptsLegacyPayloadShape(t *testing.T) {
	h := NewRecordHandler(usecase.NewRecordService(testRepo{}))
	body := []byte(`{"emision":"2026-02-17 09:41:45","nave":"NYK DENEB","viaje":"072E","cliente":"CAPITAL PACIFICO, S.A.","booking":"YMLUL160382911","contenedor":"YMLU5374938","puerto_descargue":"RODMAN","libre_retencion_hasta":"2021-03-06","dias_libre":0,"transportista":"GLOBERUNNERS, INC","titulo_terminal":"PANAMA PORTS COMPANY (RODMAN)","usuario_firma":"Admin"}`)
	r := httptest.NewRequest(http.MethodPost, "/v1/records", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(middleware.WithClaims(r.Context(), &auth.Claims{Subject: "user-1"}))

	w := httptest.NewRecorder()
	h.Create(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestValidateRecordHandler(t *testing.T) {
	secret := "test-qr-secret"
	verifier := usecase.NewCompactQRTokenVerifier(secret)
	h := NewRecordHandler(usecase.NewRecordService(testRepo{}, verifier))

	token := signedCompactToken(123, secret, time.Now().Add(10*time.Minute).Unix())
	r := httptest.NewRequest(http.MethodGet, "/v1/records/validate?t="+token, nil)
	r = r.WithContext(middleware.WithClaims(r.Context(), &auth.Claims{Subject: "user-1"}))

	w := httptest.NewRecorder()
	h.Validate(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func signedCompactToken(recordID int64, secret string, exp int64) string {
	body := fmt.Sprintf("v1|%d|%d", recordID, exp)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	sig := mac.Sum(nil)[:16]
	return fmt.Sprintf("v1.%d.%d.%s", recordID, exp, base64.RawURLEncoding.EncodeToString(sig))
}
