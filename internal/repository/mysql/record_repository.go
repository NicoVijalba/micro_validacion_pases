package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/example/validacion-pases/internal/domain"
	"github.com/go-sql-driver/mysql"
)

type RecordRepository struct {
	db *sql.DB
}

func NewRecordRepository(db *sql.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

func (r *RecordRepository) Insert(ctx context.Context, record domain.Record) (int64, error) {
	const q = `
INSERT INTO records (
    emision, nave, viaje, cliente, booking, rama, contenedor,
    puerto_descargue, libre_retencion_hasta, dias_libre, transportista,
    titulo_terminal, usuario_firma, created_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, q,
		record.Emision,
		record.Nave,
		record.Viaje,
		record.Cliente,
		record.Booking,
		record.Rama,
		record.Contenedor,
		record.PuertoDescargue,
		record.LibreRetencionHasta,
		record.DiasLibre,
		record.Transportista,
		record.TituloTerminal,
		record.UsuarioFirma,
		record.CreatedAt,
	)
	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return 0, domain.ErrConflict
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RecordRepository) FindByID(ctx context.Context, id int64) (domain.Record, error) {
	const q = `
SELECT id, emision, nave, viaje, cliente, booking, rama, contenedor, puerto_descargue,
       libre_retencion_hasta, dias_libre, transportista, titulo_terminal, usuario_firma, created_at
FROM records
WHERE id = ?`

	var rec domain.Record
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&rec.ID,
		&rec.Emision,
		&rec.Nave,
		&rec.Viaje,
		&rec.Cliente,
		&rec.Booking,
		&rec.Rama,
		&rec.Contenedor,
		&rec.PuertoDescargue,
		&rec.LibreRetencionHasta,
		&rec.DiasLibre,
		&rec.Transportista,
		&rec.TituloTerminal,
		&rec.UsuarioFirma,
		&rec.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Record{}, domain.ErrNotFound
		}
		return domain.Record{}, err
	}

	return rec, nil
}
