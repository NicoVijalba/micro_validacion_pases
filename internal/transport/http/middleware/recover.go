package middleware

import (
	"fmt"
	"net/http"

	"github.com/example/validacion-pases/pkg/problem"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				problem.Write(w, r, problem.Internal(fmt.Sprintf("panic recovered: %v", rec)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
