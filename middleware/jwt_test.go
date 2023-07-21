package middleware

import (
	"net/http"
	"testing"
)

func TestGetPhoneFromHeader(t *testing.T) {
	m := &JWTMiddleware{}

	req1 := http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer token1"},
		},
	}

	req2 := http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer token2"},
		},
	}

	phones := m.GetPhoneFromHeader(&req1, &req2)

	if len(phones) != 2 {
		t.Fatalf("Expected 2 phone numbers, got %v", len(phones))
	}

	if phones[0] != "Bearer token1" {
		t.Errorf("Expected 'Bearer token1', got %v", phones[0])
	}

	if phones[1] != "Bearer token2" {
		t.Errorf("Expected 'Bearer token2', got %v", phones[1])
	}
}
