package engine

import (
	"testing"
	"time"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"github.com/zhangkui/go-alert-evaluator/internal/silence"
	"github.com/zhangkui/go-alert-evaluator/internal/store"
)

func TestFiringWaitsForRecoveryDuration(t *testing.T) {
	start := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	series := store.NewSeriesStore(time.Hour)
	labels := model.Labels{"service": "billing"}
	engine := New(series, silence.NewStore(), nil)
	if err := engine.AddRule(model.Rule{ID: "latency", Metric: "latency", Labels: labels, Aggregation: model.AggregationMaximum, Comparator: model.ComparatorAbove, Threshold: 50, Window: 10*time.Second, RecoverFor: 30*time.Second}); err != nil { t.Fatal(err) }
	if err := series.Write("latency", labels, model.Sample{Timestamp: start, Value: 100}, start); err != nil { t.Fatal(err) }
	firing, err := engine.Evaluate("latency", start)
	if err != nil || firing.State != model.StateFiring { t.Fatalf("expected firing: %#v, %v", firing, err) }
	firstHealthy := start.Add(20*time.Second)
	if err := series.Write("latency", labels, model.Sample{Timestamp: firstHealthy, Value: 10}, firstHealthy); err != nil { t.Fatal(err) }
	recovering, err := engine.Evaluate("latency", firstHealthy)
	if err != nil { t.Fatal(err) }
	if recovering.State != model.StateFiring { t.Fatalf("state = %s, want firing until recover_for elapses", recovering.State) }
	fullyHealthy := start.Add(51*time.Second)
	if err := series.Write("latency", labels, model.Sample{Timestamp: fullyHealthy, Value: 10}, fullyHealthy); err != nil { t.Fatal(err) }
	recovered, err := engine.Evaluate("latency", fullyHealthy)
	if err != nil || recovered.State != model.StateNormal { t.Fatalf("expected normal after recovery duration: %#v, %v", recovered, err) }
}
