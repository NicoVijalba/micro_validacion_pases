package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenIssuer struct {
	alg      string
	issuer   string
	audience string
	ttl      time.Duration
	hsSecret []byte
}

func NewTokenIssuer(alg, issuer, audience string, ttl time.Duration, hsSecret string) (*TokenIssuer, error) {
	if alg != "HS256" {
		return nil, errors.New("token issuance requires JWT_ALG=HS256")
	}
	if hsSecret == "" {
		return nil, errors.New("JWT_HS_SECRET is required for token issuance")
	}

	return &TokenIssuer{
		alg:      alg,
		issuer:   issuer,
		audience: audience,
		ttl:      ttl,
		hsSecret: []byte(hsSecret),
	}, nil
}

func (i *TokenIssuer) Issue(subject string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(i.ttl)

	claims := Claims{
		Subject: subject,
		Scopes:  []string{"records:write"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings{i.audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(i.hsSecret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}
