package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"github.com/zhangkui/go-alert-evaluator/internal/clock"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
)

type cancelSender struct { cancel context.CancelFunc; calls int }
func (s *cancelSender) Send(model.Rule, model.Evaluation) error { s.calls++; if s.calls == 1 { s.cancel() }; return nil }

func TestEvaluateBatchStopsAfterContextCancellation(t *testing.T) {
	now := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	ctx, cancel := context.WithCancel(context.Background())
	sender := &cancelSender{cancel: cancel}
	service := New(clock.Fixed{Time: now}, sender, time.Hour)
	labels := model.Labels{"host": "api-1"}
	if err := service.WriteSample("cpu", labels, model.Sample{Timestamp: now, Value: 90}); err != nil { t.Fatal(err) }
	for _, id := range []string{"first", "second"} {
		if err := service.AddRule(model.Rule{ID: id, Metric: "cpu", Labels: labels, Aggregation: model.AggregationMaximum, Comparator: model.ComparatorAbove, Threshold: 50, Window: time.Minute}); err != nil { t.Fatal(err) }
	}
	results, err := service.EvaluateBatch(ctx, []string{"first", "second"}, now)
	if !errors.Is(err, context.Canceled) { t.Fatalf("error = %v, want context canceled; results=%d", err, len(results)) }
	if sender.calls != 1 { t.Fatalf("sender calls = %d, want 1", sender.calls) }
}
