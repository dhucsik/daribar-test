package middleware

import (
	"net/http"
)

type JWTMiddleware struct {
}

func (m *JWTMiddleware) GetPhoneFromHeader(reqs ...*http.Request) []string {
	out := make([]string, len(reqs))
	for _, req := range reqs {
		authHeader := req.Header.Get("Authorization")
		out = append(out, authHeader)
	}

	return out
}
