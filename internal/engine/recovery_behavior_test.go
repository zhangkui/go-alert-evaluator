package engine

import (
	"testing"
	"time"

	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"github.com/zhangkui/go-alert-evaluator/internal/silence"
	"github.com/zhangkui/go-alert-evaluator/internal/store"
)

func TestRecoveryDurationResetsAfterAnotherBreach(t *testing.T) {
	start := time.Date(2026, 8, 17, 9, 0, 0, 0, time.UTC)
	engine, series, labels := newRecoveryTestEngine(t, 30*time.Second, 0, silence.NewStore())
	assertEvaluationState(t, engine, series, labels, start, 100, model.StateFiring)
	assertEvaluationState(t, engine, series, labels, start.Add(10*time.Second), 10, model.StateFiring)
	assertEvaluationState(t, engine, series, labels, start.Add(25*time.Second), 100, model.StateFiring)
	assertEvaluationState(t, engine, series, labels, start.Add(30*time.Second), 10, model.StateFiring)
	assertEvaluationState(t, engine, series, labels, start.Add(59*time.Second), 10, model.StateFiring)
	assertEvaluationState(t, engine, series, labels, start.Add(60*time.Second), 10, model.StateNormal)
}

func TestRecoveryDurationZeroRecoversImmediately(t *testing.T) {
	start := time.Date(2026, 8, 17, 10, 0, 0, 0, time.UTC)
	engine, series, labels := newRecoveryTestEngine(t, 0, 0, silence.NewStore())
	assertEvaluationState(t, engine, series, labels, start, 100, model.StateFiring)
	assertEvaluationState(t, engine, series, labels, start.Add(time.Second), 10, model.StateNormal)
}

func TestPendingReturnsToNormalWhenHealthy(t *testing.T) {
	start := time.Date(2026, 8, 17, 11, 0, 0, 0, time.UTC)
	engine, series, labels := newRecoveryTestEngine(t, 30*time.Second, time.Minute, silence.NewStore())
	assertEvaluationState(t, engine, series, labels, start, 100, model.StatePending)
	assertEvaluationState(t, engine, series, labels, start.Add(time.Second), 10, model.StateNormal)
}

func TestInsufficientDataKeepsStateAndReasonExplicit(t *testing.T) {
	start := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
	engine, _, _ := newRecoveryTestEngine(t, 30*time.Second, 0, silence.NewStore())
	evaluation, err := engine.Evaluate("latency", start)
	if err != nil {
		t.Fatal(err)
	}
	if evaluation.State != model.StateInsufficient {
		t.Fatalf("state = %s, want %s", evaluation.State, model.StateInsufficient)
	}
	if evaluation.MissingReason == "" {
		t.Fatal("expected a missing reason")
	}
}

func TestSilenceOnlySuppressesNotificationsDuringRecovery(t *testing.T) {
	start := time.Date(2026, 8, 17, 13, 0, 0, 0, time.UTC)
	labels := model.Labels{"service": "billing"}
	silences := silence.NewStore()
	if err := silences.Add(model.Silence{ID: "maintenance", Matchers: labels, StartsAt: start, EndsAt: start.Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	engine, series, labels := newRecoveryTestEngine(t, 30*time.Second, 0, silences)
	firing := assertEvaluationState(t, engine, series, labels, start, 100, model.StateFiring)
	if firing.Notification != "suppressed" {
		t.Fatalf("notification = %q, want suppressed", firing.Notification)
	}
	recovering := assertEvaluationState(t, engine, series, labels, start.Add(10*time.Second), 10, model.StateFiring)
	if recovering.Notification != "none" {
		t.Fatalf("notification = %q, want none", recovering.Notification)
	}
	recovered := assertEvaluationState(t, engine, series, labels, start.Add(40*time.Second), 10, model.StateNormal)
	if recovered.Notification != "suppressed" {
		t.Fatalf("notification = %q, want suppressed", recovered.Notification)
	}
}

func newRecoveryTestEngine(t *testing.T, recoverFor, firingFor time.Duration, silences *silence.Store) (*Engine, *store.SeriesStore, model.Labels) {
	t.Helper()
	series := store.NewSeriesStore(time.Hour)
	labels := model.Labels{"service": "billing"}
	engine := New(series, silences, nil)
	err := engine.AddRule(model.Rule{ID: "latency", Metric: "latency", Labels: labels, Aggregation: model.AggregationMaximum, Comparator: model.ComparatorAbove, Threshold: 50, Window: time.Second, For: firingFor, RecoverFor: recoverFor})
	if err != nil {
		t.Fatal(err)
	}
	return engine, series, labels
}

func assertEvaluationState(t *testing.T, engine *Engine, series *store.SeriesStore, labels model.Labels, at time.Time, value float64, want model.AlertState) model.Evaluation {
	t.Helper()
	if err := series.Write("latency", labels, model.Sample{Timestamp: at, Value: value}, at); err != nil {
		t.Fatal(err)
	}
	evaluation, err := engine.Evaluate("latency", at)
	if err != nil {
		t.Fatal(err)
	}
	if evaluation.State != want {
		t.Fatalf("state at %s = %s, want %s", at, evaluation.State, want)
	}
	return evaluation
}
