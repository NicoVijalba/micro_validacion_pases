package problem

import (
	"encoding/json"
	"net/http"
)

type Details struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

func Write(w http.ResponseWriter, r *http.Request, p Details) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	p.Instance = r.URL.Path
	_ = json.NewEncoder(w).Encode(p)
}

func BadRequest(detail string) Details {
	return Details{Type: "about:blank", Title: "Bad Request", Status: http.StatusBadRequest, Detail: detail}
}
func Unauthorized(detail string) Details {
	return Details{Type: "about:blank", Title: "Unauthorized", Status: http.StatusUnauthorized, Detail: detail}
}
func Forbidden(detail string) Details {
	return Details{Type: "about:blank", Title: "Forbidden", Status: http.StatusForbidden, Detail: detail}
}
func Conflict(detail string) Details {
	return Details{Type: "about:blank", Title: "Conflict", Status: http.StatusConflict, Detail: detail}
}
func NotFound(detail string) Details {
	return Details{Type: "about:blank", Title: "Not Found", Status: http.StatusNotFound, Detail: detail}
}
func Internal(detail string) Details {
	return Details{Type: "about:blank", Title: "Internal Server Error", Status: http.StatusInternalServerError, Detail: detail}
}
func ServiceUnavailable(detail string) Details {
	return Details{Type: "about:blank", Title: "Service Unavailable", Status: http.StatusServiceUnavailable, Detail: detail}
}
