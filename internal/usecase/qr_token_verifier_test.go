package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCompactQRTokenVerifierSuccess(t *testing.T) {
	secret := "my-qr-secret"
	v := NewCompactQRTokenVerifier(secret)
	v.nowFn = func() time.Time {
		return time.Unix(1700000000, 0).UTC()
	}

	exp := time.Unix(1700000600, 0).UTC().Unix()
	token := signCompactToken(45, secret, exp)

	id, err := v.VerifyAndExtractRecordID(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 45 {
		t.Fatalf("expected id 45, got %d", id)
	}
}

func TestCompactQRTokenVerifierExpired(t *testing.T) {
	secret := "my-qr-secret"
	v := NewCompactQRTokenVerifier(secret)
	v.nowFn = func() time.Time {
		return time.Unix(1700000601, 0).UTC()
	}

	exp := time.Unix(1700000600, 0).UTC().Unix()
	token := signCompactToken(45, secret, exp)

	_, err := v.VerifyAndExtractRecordID(token)
	if !errors.Is(err, ErrExpiredQRToken) {
		t.Fatalf("expected ErrExpiredQRToken, got %v", err)
	}
}

func TestCompactQRTokenVerifierInvalidSignature(t *testing.T) {
	v := NewCompactQRTokenVerifier("my-qr-secret")
	_, err := v.VerifyAndExtractRecordID("v1.45.1700000600.invalidsig")
	if !errors.Is(err, ErrInvalidQRToken) {
		t.Fatalf("expected ErrInvalidQRToken, got %v", err)
	}
}

func signCompactToken(recordID int64, secret string, exp int64) string {
	body := fmt.Sprintf("v1|%d|%d", recordID, exp)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	sig := mac.Sum(nil)[:16]
	return fmt.Sprintf("v1.%d.%d.%s", recordID, exp, base64.RawURLEncoding.EncodeToString(sig))
}
