package service

import (
	"context"
	"errors"
	"github.com/zhangkui/go-alert-evaluator/internal/clock"
	"github.com/zhangkui/go-alert-evaluator/internal/engine"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"github.com/zhangkui/go-alert-evaluator/internal/silence"
	"github.com/zhangkui/go-alert-evaluator/internal/store"
	"time"
)

type Service struct {
	clock    clock.Clock
	series   *store.SeriesStore
	silences *silence.Store
	engine   *engine.Engine
}

func New(currentClock clock.Clock, sender engine.Sender, retention time.Duration) *Service {
	series := store.NewSeriesStore(retention)
	silences := silence.NewStore()
	return &Service{clock: currentClock, series: series, silences: silences, engine: engine.New(series, silences, sender)}
}

func (s *Service) WriteSample(metric string, labels model.Labels, sample model.Sample) error {
	return s.series.Write(metric, labels, sample, s.clock.Now())
}
func (s *Service) AddRule(rule model.Rule) error       { return s.engine.AddRule(rule) }
func (s *Service) AddSilence(item model.Silence) error { return s.silences.Add(item) }

func (s *Service) Evaluate(ruleID string, at time.Time) (model.Evaluation, error) {
	if at.IsZero() {
		at = s.clock.Now()
	}
	return s.engine.Evaluate(ruleID, at)
}

func (s *Service) EvaluateBatch(ctx context.Context, ruleIDs []string, at time.Time) ([]model.Evaluation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	results := make([]model.Evaluation, 0, len(ruleIDs))
	for _, ruleID := range ruleIDs {
		result, err := s.Evaluate(ruleID, at)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (s *Service) History(ruleID string) ([]model.Evaluation, error) {
	history := s.engine.History(ruleID)
	if len(history) == 0 {
		return nil, errors.New("history not found")
	}
	return history, nil
}

func (s *Service) Cleanup() int { return s.series.Cleanup(s.clock.Now()) }
