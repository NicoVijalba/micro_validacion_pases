package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidQRToken        = errors.New("invalid qr token")
	ErrExpiredQRToken        = errors.New("expired qr token")
	ErrQRVerifierUnavailable = errors.New("qr verifier unavailable")
)

type QRTokenVerifier interface {
	VerifyAndExtractRecordID(token string) (int64, error)
}

type CompactQRTokenVerifier struct {
	secret []byte
	nowFn  func() time.Time
}

func NewCompactQRTokenVerifier(secret string) *CompactQRTokenVerifier {
	return &CompactQRTokenVerifier{
		secret: []byte(strings.TrimSpace(secret)),
		nowFn:  time.Now,
	}
}

func (v *CompactQRTokenVerifier) VerifyAndExtractRecordID(token string) (int64, error) {
	if len(v.secret) == 0 {
		return 0, ErrQRVerifierUnavailable
	}

	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 4 {
		return 0, ErrInvalidQRToken
	}

	version := parts[0]
	if version != "v1" {
		return 0, ErrInvalidQRToken
	}

	recordID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || recordID <= 0 {
		return 0, ErrInvalidQRToken
	}

	exp, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || exp <= 0 {
		return 0, ErrInvalidQRToken
	}

	body := fmt.Sprintf("%s|%d|%d", version, recordID, exp)
	mac := hmac.New(sha256.New, v.secret)
	_, _ = mac.Write([]byte(body))
	expected := mac.Sum(nil)[:16]
	provided, err := decodeBase64URL(parts[3])
	if err != nil {
		return 0, ErrInvalidQRToken
	}
	if !hmac.Equal(expected, provided) {
		return 0, ErrInvalidQRToken
	}

	if v.nowFn().UTC().Unix() > exp {
		return 0, ErrExpiredQRToken
	}

	return recordID, nil
}

func decodeBase64URL(raw string) ([]byte, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("empty")
	}
	padded := raw
	if m := len(raw) % 4; m != 0 {
		padded += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(padded)
}
