package engine

import (
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"github.com/zhangkui/go-alert-evaluator/internal/silence"
	"github.com/zhangkui/go-alert-evaluator/internal/store"
	"testing"
	"time"
)

func TestEvaluateTransitionsPendingToFiring(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	series := store.NewSeriesStore(time.Hour)
	labels := model.Labels{"service": "api"}
	if err := series.Write("latency", labels, model.Sample{Timestamp: now, Value: 100}, now); err != nil {
		t.Fatal(err)
	}
	engine := New(series, silence.NewStore(), nil)
	if err := engine.AddRule(model.Rule{ID: "r1", Metric: "latency", Labels: labels, Aggregation: model.AggregationMaximum, Comparator: model.ComparatorAbove, Threshold: 50, Window: time.Minute, For: 30 * time.Second}); err != nil {
		t.Fatal(err)
	}
	first, err := engine.Evaluate("r1", now)
	if err != nil || first.State != model.StatePending {
		t.Fatalf("expected pending: %#v, %v", first, err)
	}
	if err := series.Write("latency", labels, model.Sample{Timestamp: now.Add(31 * time.Second), Value: 100}, now.Add(31*time.Second)); err != nil {
		t.Fatal(err)
	}
	second, err := engine.Evaluate("r1", now.Add(31*time.Second))
	if err != nil || second.State != model.StateFiring {
		t.Fatalf("expected firing: %#v, %v", second, err)
	}
}

func TestPartialSilenceMatchDoesNotSuppressNotificationOrHistory(t *testing.T) {
	now := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	series := store.NewSeriesStore(time.Hour)
	labels := model.Labels{"service": "api", "region": "east"}
	if err := series.Write("latency", labels, model.Sample{Timestamp: now, Value: 100}, now); err != nil {
		t.Fatal(err)
	}
	silences := silence.NewStore()
	if err := silences.Add(model.Silence{ID: "west-api", Matchers: model.Labels{"service": "api", "region": "west"}, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}); err != nil {
		t.Fatal(err)
	}
	engine := New(series, silences, nil)
	if err := engine.AddRule(model.Rule{ID: "r1", Metric: "latency", Labels: labels, Aggregation: model.AggregationMaximum, Comparator: model.ComparatorAbove, Threshold: 50, Window: time.Minute}); err != nil {
		t.Fatal(err)
	}
	evaluation, err := engine.Evaluate("r1", now)
	if err != nil {
		t.Fatal(err)
	}
	if evaluation.Notification != "firing" || !evaluation.NotificationSent {
		t.Fatalf("expected notification to be sent, got %#v", evaluation)
	}
	history := engine.History("r1")
	if len(history) != 1 || history[0].Notification != "firing" || history[0].NotificationSent != true {
		t.Fatalf("unexpected history: %#v", history)
	}
}
