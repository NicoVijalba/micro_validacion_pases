package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Subject string
	Scopes  []string `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

func (c Claims) HasScope(want string) bool {
	for _, s := range c.Scopes {
		if s == want {
			return true
		}
	}
	return false
}

type JWTValidator struct {
	alg       string
	issuer    string
	audience  string
	clockSkew time.Duration
	hsSecret  []byte
}

func NewJWTValidator(_ context.Context, alg, issuer, audience string, clockSkew time.Duration, hsSecret, jwksURL string, _ time.Duration) (*JWTValidator, error) {
	v := &JWTValidator{
		alg:       alg,
		issuer:    issuer,
		audience:  audience,
		clockSkew: clockSkew,
		hsSecret:  []byte(hsSecret),
	}

	switch alg {
	case "HS256":
		if hsSecret == "" {
			return nil, errors.New("JWT_HS_SECRET is required for HS256")
		}
		return v, nil
	case "RS256":
		if jwksURL == "" {
			return nil, errors.New("JWT_JWKS_URL is required for RS256")
		}
		return nil, errors.New("RS256 validation is temporarily disabled in this local build; use HS256")
	default:
		return nil, errors.New("unsupported jwt algorithm")
	}
}

func (v *JWTValidator) Parse(token string) (*Claims, error) {
	if v.alg != "HS256" {
		return nil, errors.New("only HS256 validation is enabled")
	}

	claims := &Claims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{v.alg}),
		jwt.WithAudience(v.audience),
		jwt.WithIssuer(v.issuer),
		jwt.WithLeeway(v.clockSkew),
	)

	_, err := parser.ParseWithClaims(token, claims, func(_ *jwt.Token) (any, error) {
		return v.hsSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims.Subject == "" {
		claims.Subject = claims.RegisteredClaims.Subject
	}
	if claims.Subject == "" {
		return nil, errors.New("token missing subject")
	}

	return claims, nil
}
