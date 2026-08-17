package service

import (
	"context"
	"github.com/zhangkui/go-alert-evaluator/internal/clock"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"testing"
	"time"
)

func TestEvaluateBatch(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	service := New(clock.Fixed{Time: now}, nil, time.Hour)
	labels := model.Labels{"host": "one"}
	if err := service.WriteSample("cpu", labels, model.Sample{Timestamp: now, Value: 80}); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"r1", "r2"} {
		if err := service.AddRule(model.Rule{ID: id, Metric: "cpu", Labels: labels, Aggregation: model.AggregationMaximum, Comparator: model.ComparatorAbove, Threshold: 70, Window: time.Minute}); err != nil {
			t.Fatal(err)
		}
	}
	results, err := service.EvaluateBatch(context.Background(), []string{"r1", "r2"}, now)
	if err != nil || len(results) != 2 {
		t.Fatalf("unexpected batch result: %#v, %v", results, err)
	}
}
