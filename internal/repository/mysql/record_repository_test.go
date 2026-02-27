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
	defer func() {
		if cerr := db.Close(); cerr != nil {
			t.Errorf("failed to close db: %v", cerr)
		}
	}()

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
	mock.ExpectClose()

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

func TestFindByIDSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			t.Errorf("failed to close db: %v", cerr)
		}
	}()

	repo := NewRecordRepository(db)
	now := time.Date(2026, 2, 17, 9, 41, 45, 0, time.UTC)
	lrh := time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{
		"id", "emision", "nave", "viaje", "cliente", "booking", "rama", "contenedor", "puerto_descargue",
		"libre_retencion_hasta", "dias_libre", "transportista", "titulo_terminal", "usuario_firma", "created_at",
	}).AddRow(
		int64(10), now, "NYK DENEB", "072E", "CAPITAL PACIFICO, S.A.", "YMLUL160382911", "internacional", "YMLU5374938", "RODMAN",
		lrh, 17, "", "PANAMA PORTS COMPANY (RODMAN)", "Admin", now,
	)

	mock.ExpectQuery("SELECT id, emision, nave").WithArgs(int64(10)).WillReturnRows(rows)
	mock.ExpectClose()

	rec, err := repo.FindByID(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.ID != 10 {
		t.Fatalf("expected id 10, got %d", rec.ID)
	}
}
