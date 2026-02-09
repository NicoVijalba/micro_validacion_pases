package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/example/validacion-pases/internal/usecase"
	"github.com/example/validacion-pases/pkg/problem"
)

type TokenHandler struct {
	service  *usecase.TokenService
	validate *validator.Validate
}

type tokenRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=200"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
}

func NewTokenHandler(service *usecase.TokenService) *TokenHandler {
	return &TokenHandler{service: service, validate: validator.New()}
}

func (h *TokenHandler) Issue(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
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

	token, expiresAt, err := h.service.Issue(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			problem.Write(w, r, problem.Unauthorized("invalid credentials"))
			return
		}
		problem.Write(w, r, problem.ServiceUnavailable("token service unavailable"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
	})
}
