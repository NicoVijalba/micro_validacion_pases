package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/example/validacion-pases/internal/security/auth"
	"github.com/example/validacion-pases/pkg/problem"
)

type contextKey string

const claimsContextKey contextKey = "claims"

func AuthBearer(validator *auth.JWTValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				problem.Write(w, r, problem.Unauthorized("missing or invalid bearer token"))
				return
			}

			raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			claims, err := validator.Parse(raw)
			if err != nil {
				problem.Write(w, r, problem.Unauthorized("invalid token"))
				return
			}

			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
		})
	}
}

func WithClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*auth.Claims, error) {
	claims, ok := ctx.Value(claimsContextKey).(*auth.Claims)
	if !ok || claims == nil {
		return nil, errors.New("claims not found")
	}
	return claims, nil
}
