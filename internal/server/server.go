package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"spc-chart/internal/capability"
	"spc-chart/internal/chart"
)

type Config struct {
	Addr string
}

func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/ichart", handleIChart)
	mux.HandleFunc("/api/capability", handleCapability)
	mux.HandleFunc("/api/cusum", handleCUSUM)
	return mux
}

func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type ichartRequest struct {
	Values []float64 `json:"values"`
	Sigma  float64   `json:"sigma"`
}

type limitsOutput struct {
	CL  float64 `json:"cl"`
	UCL float64 `json:"ucl"`
	LCL float64 `json:"lcl"`
}

type ichartResponse struct {
	Type     string       `json:"type"`
	Limits   limitsOutput `json:"limits"`
	OOCCount int          `json:"ooc_count"`
	Mean     float64      `json:"mean"`
}

func handleIChart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req ichartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Values) < 2 {
		httpError(w, http.StatusBadRequest, "at least 2 values required")
		return
	}
	if req.Sigma <= 0 {
		req.Sigma = 3
	}
	res, err := chart.Individuals(req.Values, chart.IndividualsConfig{Sigma: req.Sigma, MRSpan: 2})
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	out := ichartResponse{
		Type:     res.Type.String(),
		Limits:   limitsOutput{CL: res.Limits.CL, UCL: res.Limits.UCL, LCL: res.Limits.LCL},
		OOCCount: res.OOCCount,
		Mean:     res.Mean,
	}
	holdIChartUCL(&out)
	writeJSON(w, http.StatusOK, out)
}

type capabilityRequest struct {
	Values []float64 `json:"values"`
	USL    float64   `json:"usl"`
	LSL    float64   `json:"lsl"`
	Target float64   `json:"target"`
}

type capabilityResponse struct {
	N    int     `json:"n"`
	Mean float64 `json:"mean"`
	Cp   float64 `json:"cp"`
	Cpk  float64 `json:"cpk"`
	Pp   float64 `json:"pp"`
	Ppk  float64 `json:"ppk"`
	Cpm  float64 `json:"cpm"`
}

func handleCapability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req capabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Values) < 2 {
		httpError(w, http.StatusBadRequest, "at least 2 values required")
		return
	}
	spec := capability.SpecLimits{USL: req.USL, LSL: req.LSL, Target: req.Target}
	res, err := capability.ComputeIndices(req.Values, spec, 0)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, capabilityResponse{
		N: res.N, Mean: res.Mean, Cp: res.Cp, Cpk: res.Cpk,
		Pp: res.Pp, Ppk: res.Ppk, Cpm: res.Cpm,
	})
}

type cusumRequest struct {
	Values []float64 `json:"values"`
	K      float64   `json:"k"`
	H      float64   `json:"h"`
	Target float64   `json:"target"`
}

type cusumResponse struct {
	OOCCount int     `json:"ooc_count"`
	Target   float64 `json:"target"`
	Sigma    float64 `json:"sigma"`
}

func handleCUSUM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req cusumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Values) < 2 {
		httpError(w, http.StatusBadRequest, "at least 2 values required")
		return
	}
	if req.K <= 0 {
		req.K = 0.5
	}
	if req.H <= 0 {
		req.H = 5.0
	}
	cfg := chart.CUSUMConfig{Target: req.Target, K: req.K, H: req.H}
	res, err := chart.CUSUM(req.Values, cfg)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cusumResponse{
		OOCCount: res.OOCCount,
		Target:   res.Target,
		Sigma:    res.Sigma,
	})
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
