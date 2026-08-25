package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestIChartEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := ichartRequest{Values: []float64{10, 11, 9, 10, 12, 10, 11, 9, 10, 10}, Sigma: 3}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/ichart", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp ichartResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Limits.UCL <= resp.Limits.CL {
		t.Error("expected UCL > CL")
	}
}

func TestIChartEndpoint_TooFew(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	body := []byte(`{"values":[1]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/ichart", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCapabilityEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := capabilityRequest{
		Values: []float64{10, 10.1, 9.9, 10.2, 9.8, 10.1, 10, 9.9, 10.1, 10},
		USL:    11, LSL: 9, Target: 10,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/capability", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp capabilityResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Cp <= 0 {
		t.Error("expected Cp > 0")
	}
}

func TestCUSUMEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := cusumRequest{Values: []float64{0, 0.1, -0.1, 0.2, 0, 0.1, 5, 5.1, 5, 5.2}, K: 0.5, H: 5}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/cusum", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	endpoints := []string{"/api/ichart", "/api/capability", "/api/cusum"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", ep, rec.Code)
		}
	}
}
