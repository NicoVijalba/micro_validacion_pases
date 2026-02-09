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
