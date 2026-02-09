package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/example/validacion-pases/internal/domain"
)

type RecordService struct {
	repo domain.RecordRepository
}

func NewRecordService(repo domain.RecordRepository) *RecordService {
	return &RecordService{repo: repo}
}

func (s *RecordService) Create(ctx context.Context, in domain.CreateRecordInput) (int64, domain.Record, error) {
	if strings.TrimSpace(in.UsuarioFirma) == "" {
		return 0, domain.Record{}, domain.ErrUnauthorized
	}
	if strings.TrimSpace(in.Nave) == "" || strings.TrimSpace(in.Viaje) == "" || strings.TrimSpace(in.Cliente) == "" {
		return 0, domain.Record{}, domain.ErrInvalidInput
	}
	if strings.TrimSpace(in.Booking) == "" || strings.TrimSpace(in.PuertoDescargue) == "" {
		return 0, domain.Record{}, domain.ErrInvalidInput
	}

	diasLibre := 0
	if in.DiasLibre != nil {
		if *in.DiasLibre < 0 {
			return 0, domain.Record{}, domain.ErrInvalidInput
		}
		diasLibre = *in.DiasLibre
	}

	contenedor, transportista, err := resolveContenedorData(in.Rama, in.ContenedorSerie, in.CodigoISO, in.Transportista)
	if err != nil {
		return 0, domain.Record{}, domain.ErrInvalidInput
	}

	rec := domain.Record{
		Emision:             time.Now().UTC(),
		Nave:                strings.TrimSpace(in.Nave),
		Viaje:               strings.TrimSpace(in.Viaje),
		Cliente:             strings.TrimSpace(in.Cliente),
		Booking:             strings.TrimSpace(in.Booking),
		Rama:                strings.ToLower(strings.TrimSpace(in.Rama)),
		Contenedor:          contenedor,
		PuertoDescargue:     strings.TrimSpace(in.PuertoDescargue),
		LibreRetencionHasta: in.FechaReal.AddDate(0, 0, diasLibre),
		DiasLibre:           diasLibre,
		Transportista:       transportista,
		TituloTerminal:      resolveTituloTerminal(in.PuertoDescargue),
		UsuarioFirma:        strings.TrimSpace(in.UsuarioFirma),
		CreatedAt:           time.Now().UTC(),
	}

	id, err := s.repo.Insert(ctx, rec)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return 0, domain.Record{}, domain.ErrConflict
		}
		return 0, domain.Record{}, err
	}

	rec.ID = id
	return id, rec, nil
}

func resolveContenedorData(rama, contenedorSerie, codigoISO, transportista string) (string, string, error) {
	normalizedRama := strings.ToLower(strings.TrimSpace(rama))
	switch normalizedRama {
	case "internacional":
		serie := strings.TrimSpace(contenedorSerie)
		if serie == "" {
			return "", "", errors.New("contenedor_serie is required for internacional")
		}
		return serie, "", nil
	case "nacional":
		iso := strings.TrimSpace(codigoISO)
		trans := strings.TrimSpace(transportista)
		if iso == "" || trans == "" {
			return "", "", errors.New("codigo_iso and transportista are required for nacional")
		}
		return fmt.Sprintf("1 X %s", iso), trans, nil
	default:
		return "", "", errors.New("rama must be internacional or nacional")
	}
}

func resolveTituloTerminal(puertoDescargue string) string {
	puerto := strings.ToUpper(strings.TrimSpace(puertoDescargue))
	switch {
	case strings.Contains(puerto, "BALBOA"):
		return "TERMINAL PACIFICO - BALBOA"
	case strings.Contains(puerto, "CRISTOBAL"):
		return "TERMINAL ATLANTICO - CRISTOBAL"
	case strings.Contains(puerto, "MANZANILLO"):
		return "MANZANILLO INTERNATIONAL TERMINAL"
	default:
		return "TERMINAL " + strings.TrimSpace(puertoDescargue)
	}
}
