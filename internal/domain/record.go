package domain

import (
	"context"
	"time"
)

// Record represents a persisted pass validation record.
type Record struct {
	ID                  int64
	Emision             time.Time
	Nave                string
	Viaje               string
	Cliente             string
	Booking             string
	Rama                string
	Contenedor          string
	PuertoDescargue     string
	LibreRetencionHasta time.Time
	DiasLibre           int
	Transportista       string
	TituloTerminal      string
	UsuarioFirma        string
	CreatedAt           time.Time
}

// CreateRecordInput contains the required fields to create a new record.
type CreateRecordInput struct {
	Nave            string
	Viaje           string
	Cliente         string
	Booking         string
	Rama            string
	ContenedorSerie string
	CodigoISO       string
	FechaReal       time.Time
	DiasLibre       *int
	Transportista   string
	PuertoDescargue string
	UsuarioFirma    string
}

// RecordRepository defines persistence operations for records.
type RecordRepository interface {
	Insert(ctx context.Context, record Record) (int64, error)
	FindByID(ctx context.Context, id int64) (Record, error)
}
