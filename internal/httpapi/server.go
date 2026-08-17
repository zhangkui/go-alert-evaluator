package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"github.com/zhangkui/go-alert-evaluator/internal/service"
	"net/http"
	"strings"
	"time"
)

type Server struct{ service *service.Service }

func New(app *service.Service) http.Handler {
	server := &Server{service: app}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /samples", server.writeSample)
	mux.HandleFunc("POST /rules", server.addRule)
	mux.HandleFunc("POST /evaluate", server.evaluate)
	mux.HandleFunc("POST /evaluate/batch", server.evaluateBatch)
	mux.HandleFunc("POST /silences", server.addSilence)
	mux.HandleFunc("GET /alerts/{ruleID}/history", server.history)
	return mux
}

type sampleRequest struct {
	Metric    string       `json:"metric"`
	Labels    model.Labels `json:"labels"`
	Timestamp time.Time    `json:"timestamp"`
	Value     float64      `json:"value"`
}

func (s *Server) writeSample(w http.ResponseWriter, r *http.Request) {
	var request sampleRequest
	if !decode(w, r, &request) {
		return
	}
	if err := s.service.WriteSample(request.Metric, request.Labels, model.Sample{Timestamp: request.Timestamp, Value: request.Value}); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type ruleRequest struct {
	ID          string            `json:"id"`
	Metric      string            `json:"metric"`
	Labels      model.Labels      `json:"labels"`
	Aggregation model.Aggregation `json:"aggregation"`
	Comparator  model.Comparator  `json:"comparator"`
	Threshold   float64           `json:"threshold"`
	Window      string            `json:"window"`
	For         string            `json:"for"`
	RecoverFor  string            `json:"recover_for"`
}

func (s *Server) addRule(w http.ResponseWriter, r *http.Request) {
	var request ruleRequest
	if !decode(w, r, &request) {
		return
	}
	window, err := time.ParseDuration(request.Window)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	forDuration, err := parseOptionalDuration(request.For)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	recoverFor, err := parseOptionalDuration(request.RecoverFor)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	rule := model.Rule{ID: request.ID, Metric: request.Metric, Labels: request.Labels, Aggregation: request.Aggregation, Comparator: request.Comparator, Threshold: request.Threshold, Window: window, For: forDuration, RecoverFor: recoverFor}
	if err := s.service.AddRule(rule); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, rule)
}

type evaluateRequest struct {
	RuleID string    `json:"rule_id"`
	At     time.Time `json:"at"`
}

func (s *Server) evaluate(w http.ResponseWriter, r *http.Request) {
	var request evaluateRequest
	if !decode(w, r, &request) {
		return
	}
	result, err := s.service.Evaluate(request.RuleID, request.At)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type batchRequest struct {
	RuleIDs []string  `json:"rule_ids"`
	At      time.Time `json:"at"`
}

func (s *Server) evaluateBatch(w http.ResponseWriter, r *http.Request) {
	var request batchRequest
	if !decode(w, r, &request) {
		return
	}
	results, err := s.service.EvaluateBatch(r.Context(), request.RuleIDs, request.At)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, r.Context().Err()) {
			status = http.StatusRequestTimeout
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func (s *Server) addSilence(w http.ResponseWriter, r *http.Request) {
	var request model.Silence
	if !decode(w, r, &request) {
		return
	}
	if err := s.service.AddSilence(request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, request)
}

func (s *Server) history(w http.ResponseWriter, r *http.Request) {
	history, err := s.service.History(strings.TrimSpace(r.PathValue("ruleID")))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, history)
}

func decode(w http.ResponseWriter, r *http.Request, value any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
func parseOptionalDuration(value string) (time.Duration, error) {
	if value == "" {
		return 0, nil
	}
	return time.ParseDuration(value)
}
