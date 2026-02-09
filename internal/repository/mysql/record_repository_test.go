package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/example/validacion-pases/internal/domain"
)

func TestInsertSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRecordRepository(db)
	now := time.Now().UTC()
	rec := domain.Record{
		Emision:             now,
		Nave:                "NAVE 1",
		Viaje:               "VJ1",
		Cliente:             "CLIENTE 1",
		Booking:             "BK1",
		Rama:                "internacional",
		Contenedor:          "ABCU1234567",
		PuertoDescargue:     "Balboa",
		LibreRetencionHasta: now,
		DiasLibre:           2,
		Transportista:       "",
		TituloTerminal:      "TERMINAL PACIFICO - BALBOA",
		UsuarioFirma:        "user-1",
		CreatedAt:           now,
	}

	mock.ExpectExec("INSERT INTO records").WithArgs(
		rec.Emision,
		rec.Nave,
		rec.Viaje,
		rec.Cliente,
		rec.Booking,
		rec.Rama,
		rec.Contenedor,
		rec.PuertoDescargue,
		rec.LibreRetencionHasta,
		rec.DiasLibre,
		rec.Transportista,
		rec.TituloTerminal,
		rec.UsuarioFirma,
		rec.CreatedAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.Insert(context.Background(), rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewRecordRepository(t *testing.T) {
	var db sql.DB
	r := NewRecordRepository(&db)
	if r == nil {
		t.Fatal("repo nil")
	}
}
