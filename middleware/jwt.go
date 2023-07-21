package middleware

import "net/http"

type JWTMiddleware struct {
}

func (m *JWTMiddleware) GetPhoneFromHeader(r *http.Request) (string, error) {
	return "87088448227", nil
}
