package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/validacion-pases/internal/domain"
)

type mockRepo struct {
	insertFn func(ctx context.Context, r domain.Record) (int64, error)
}

func (m mockRepo) Insert(ctx context.Context, r domain.Record) (int64, error) {
	return m.insertFn(ctx, r)
}

func TestCreateSuccessInternacional(t *testing.T) {
	svc := NewRecordService(mockRepo{insertFn: func(_ context.Context, _ domain.Record) (int64, error) {
		return 99, nil
	}})

	fechaReal, _ := time.Parse("2006-01-02", "2026-02-09")
	dias := 3
	id, rec, err := svc.Create(context.Background(), domain.CreateRecordInput{
		Nave:            "NAVE TEST",
		Viaje:           "VJ001",
		Cliente:         "CLIENTE TEST",
		Booking:         "BK001",
		Rama:            "internacional",
		ContenedorSerie: "ABCU1234567",
		FechaReal:       fechaReal,
		DiasLibre:       &dias,
		PuertoDescargue: "Balboa",
		UsuarioFirma:    "user-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 99 {
		t.Fatalf("unexpected id: %d", id)
	}
	if rec.Contenedor != "ABCU1234567" {
		t.Fatalf("unexpected contenedor: %s", rec.Contenedor)
	}
	if got := rec.LibreRetencionHasta.Format("2006-01-02"); got != "2026-02-12" {
		t.Fatalf("unexpected libre retencion hasta: %s", got)
	}
}

func TestCreateSuccessNacional(t *testing.T) {
	svc := NewRecordService(mockRepo{insertFn: func(_ context.Context, _ domain.Record) (int64, error) {
		return 1, nil
	}})

	fechaReal, _ := time.Parse("2006-01-02", "2026-02-09")
	id, rec, err := svc.Create(context.Background(), domain.CreateRecordInput{
		Nave:            "NAVE TEST",
		Viaje:           "VJ001",
		Cliente:         "CLIENTE TEST",
		Booking:         "BK001",
		Rama:            "nacional",
		CodigoISO:       "22G1",
		FechaReal:       fechaReal,
		Transportista:   "TRANSPORTE SA",
		PuertoDescargue: "Cristobal",
		UsuarioFirma:    "user-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == 0 {
		t.Fatalf("unexpected id: %d", id)
	}
	if rec.Contenedor != "1 X 22G1" {
		t.Fatalf("unexpected contenedor nacional: %s", rec.Contenedor)
	}
	if rec.Transportista != "TRANSPORTE SA" {
		t.Fatalf("unexpected transportista: %s", rec.Transportista)
	}
}

func TestCreateInvalidInput(t *testing.T) {
	svc := NewRecordService(mockRepo{insertFn: func(_ context.Context, _ domain.Record) (int64, error) { return 1, nil }})
	_, _, err := svc.Create(context.Background(), domain.CreateRecordInput{Rama: "internacional", UsuarioFirma: "x"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateConflict(t *testing.T) {
	svc := NewRecordService(mockRepo{insertFn: func(_ context.Context, _ domain.Record) (int64, error) {
		return 0, domain.ErrConflict
	}})

	fechaReal, _ := time.Parse("2006-01-02", "2026-02-09")
	_, _, err := svc.Create(context.Background(), domain.CreateRecordInput{
		Nave:            "NAVE TEST",
		Viaje:           "VJ001",
		Cliente:         "CLIENTE TEST",
		Booking:         "BK001",
		Rama:            "internacional",
		ContenedorSerie: "ABCU1234567",
		FechaReal:       fechaReal,
		PuertoDescargue: "Balboa",
		UsuarioFirma:    "user-1",
	})
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
