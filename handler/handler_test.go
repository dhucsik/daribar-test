package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetSseHeaders(t *testing.T) {
	h := &Handler{}

	rw := httptest.NewRecorder()

	h.SetSseHeaders(rw)

	headers := rw.Header()

	if headers.Get("Content-Type") != "text/event-stream" {
		t.Errorf("unexpected Content-Type header: got %v want %v",
			headers.Get("Content-Type"), "text/event-stream")
	}

	if headers.Get("Connection") != "keep-alive" {
		t.Errorf("unexpected Connection header: got %v want %v",
			headers.Get("Connection"), "keep-alive")
	}

	if headers.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("unexpected Access-Control-Allow-Origin header: got %v want %v",
			headers.Get("Access-Control-Allow-Origin"), "*")
	}
}

func TestPingAndFlush(t *testing.T) {
	h := &Handler{}

	rw := httptest.NewRecorder()

	h.pingAndFlush(rw, "test data")

	result := rw.Result()

	expectedBody := "data: test data\n\n"
	body := rw.Body.String()
	if !strings.Contains(body, expectedBody) {
		t.Errorf("unexpected body: got %v want %v", body, expectedBody)
	}

	expectedStatusCode := http.StatusOK
	if result.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: got %v want %v", result.StatusCode, expectedStatusCode)
	}
}
