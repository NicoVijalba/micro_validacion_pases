package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/example/validacion-pases/internal/security/auth"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type TokenService struct {
	users  *auth.UserStore
	issuer *auth.TokenIssuer
}

func NewTokenService(users *auth.UserStore, issuer *auth.TokenIssuer) *TokenService {
	return &TokenService{users: users, issuer: issuer}
}

func (s *TokenService) Issue(username, password string) (string, time.Time, error) {
	if s.users == nil || s.issuer == nil {
		return "", time.Time{}, errors.New("token service is not configured")
	}
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return "", time.Time{}, ErrInvalidCredentials
	}
	if !s.users.Validate(username, password) {
		return "", time.Time{}, ErrInvalidCredentials
	}
	return s.issuer.Issue(strings.TrimSpace(username))
}
