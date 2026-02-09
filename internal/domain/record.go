package domain

import (
	"context"
	"time"
)

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

type RecordRepository interface {
	Insert(ctx context.Context, record Record) (int64, error)
}
