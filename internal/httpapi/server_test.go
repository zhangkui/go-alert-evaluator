package httpapi

import (
	"bytes"
	"github.com/zhangkui/go-alert-evaluator/internal/clock"
	"github.com/zhangkui/go-alert-evaluator/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWriteSampleEndpoint(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	handler := New(service.New(clock.Fixed{Time: now}, nil, time.Hour))
	body := []byte(`{"metric":"cpu","labels":{"host":"one"},"timestamp":"2026-01-01T12:00:00Z","value":42}`)
	request := httptest.NewRequest(http.MethodPost, "/samples", bytes.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("unexpected status %d: %s", response.Code, response.Body.String())
	}
}
