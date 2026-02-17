package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/example/validacion-pases/internal/domain"
	"github.com/example/validacion-pases/internal/transport/http/middleware"
	"github.com/example/validacion-pases/internal/usecase"
	"github.com/example/validacion-pases/pkg/problem"
)

type RecordHandler struct {
	service  *usecase.RecordService
	validate *validator.Validate
}

type createRecordRequest struct {
	Nave                string `json:"nave" validate:"required,max=150"`
	Viaje               string `json:"viaje" validate:"required,max=100"`
	Cliente             string `json:"cliente" validate:"required,max=200"`
	Booking             string `json:"booking" validate:"required,max=100"`
	Rama                string `json:"rama" validate:"omitempty,oneof=internacional nacional"`
	ContenedorSerie     string `json:"contenedor_serie" validate:"max=100"`
	Contenedor          string `json:"contenedor" validate:"max=100"`
	CodigoISO           string `json:"codigo_iso" validate:"max=20"`
	FechaReal           string `json:"fecha_real" validate:"omitempty,datetime=2006-01-02"`
	LibreRetencionHasta string `json:"libre_retencion_hasta" validate:"omitempty,datetime=2006-01-02"`
	DiasLibre           *int   `json:"dias_libre" validate:"omitempty,gte=0,lte=365"`
	Transportista       string `json:"transportista" validate:"max=200"`
	PuertoDescargue     string `json:"puerto_descargue" validate:"required,max=150"`
	Emision             string `json:"emision" validate:"omitempty,max=50"`
	TituloTerminal      string `json:"titulo_terminal" validate:"omitempty,max=200"`
	UsuarioFirma        string `json:"usuario_firma" validate:"omitempty,max=200"`
}

type createRecordResponse struct {
	ID                  int64  `json:"id"`
	Emision             string `json:"emision"`
	Contenedor          string `json:"contenedor"`
	LibreRetencionHasta string `json:"libre_retencion_hasta"`
	TituloTerminal      string `json:"titulo_terminal"`
	UsuarioFirma        string `json:"usuario_firma"`
}

func NewRecordHandler(service *usecase.RecordService) *RecordHandler {
	return &RecordHandler{service: service, validate: validator.New()}
}

func (h *RecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRecordRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		if errors.Is(err, io.EOF) {
			problem.Write(w, r, problem.BadRequest("empty body"))
			return
		}
		problem.Write(w, r, problem.BadRequest("invalid json payload"))
		return
	}
	if dec.More() {
		problem.Write(w, r, problem.BadRequest("multiple json values are not allowed"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		problem.Write(w, r, problem.BadRequest("payload validation failed"))
		return
	}

	if strings.TrimSpace(req.ContenedorSerie) == "" {
		req.ContenedorSerie = strings.TrimSpace(req.Contenedor)
	}

	claims, err := middleware.ClaimsFromContext(r.Context())
	if err != nil {
		problem.Write(w, r, problem.Unauthorized("auth claims missing"))
		return
	}

	var fechaReal time.Time
	switch {
	case strings.TrimSpace(req.FechaReal) != "":
		fechaReal, err = time.Parse("2006-01-02", req.FechaReal)
		if err != nil {
			problem.Write(w, r, problem.BadRequest("fecha_real must use YYYY-MM-DD"))
			return
		}
	case strings.TrimSpace(req.LibreRetencionHasta) != "":
		lrh, parseErr := time.Parse("2006-01-02", req.LibreRetencionHasta)
		if parseErr != nil {
			problem.Write(w, r, problem.BadRequest("libre_retencion_hasta must use YYYY-MM-DD"))
			return
		}
		diasLibre := 0
		if req.DiasLibre != nil {
			diasLibre = *req.DiasLibre
		}
		fechaReal = lrh.AddDate(0, 0, -diasLibre)
	default:
		problem.Write(w, r, problem.BadRequest("fecha_real is required"))
		return
	}

	id, rec, err := h.service.Create(r.Context(), domain.CreateRecordInput{
		Nave:            req.Nave,
		Viaje:           req.Viaje,
		Cliente:         req.Cliente,
		Booking:         req.Booking,
		Rama:            req.Rama,
		ContenedorSerie: req.ContenedorSerie,
		CodigoISO:       req.CodigoISO,
		FechaReal:       fechaReal,
		DiasLibre:       req.DiasLibre,
		Transportista:   req.Transportista,
		PuertoDescargue: req.PuertoDescargue,
		UsuarioFirma:    claims.Subject,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			problem.Write(w, r, problem.BadRequest("invalid input"))
		case errors.Is(err, domain.ErrConflict):
			problem.Write(w, r, problem.Conflict("record already exists"))
		case errors.Is(err, domain.ErrUnauthorized):
			problem.Write(w, r, problem.Unauthorized("unauthorized"))
		default:
			problem.Write(w, r, problem.Internal("failed to create record"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createRecordResponse{
		ID:                  id,
		Emision:             rec.Emision.Format(time.RFC3339),
		Contenedor:          rec.Contenedor,
		LibreRetencionHasta: rec.LibreRetencionHasta.Format("2006-01-02"),
		TituloTerminal:      rec.TituloTerminal,
		UsuarioFirma:        rec.UsuarioFirma,
	})
}
