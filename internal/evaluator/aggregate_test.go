package evaluator

import (
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"testing"
	"time"
)

func TestAggregateAverage(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := Aggregate(model.AggregationAverage, []model.Sample{{Timestamp: start.Add(time.Second), Value: 2}, {Timestamp: start.Add(2 * time.Second), Value: 4}}, start)
	if err != nil || result.Value != 3 {
		t.Fatalf("unexpected result: %#v, %v", result, err)
	}
}

func TestBreachedUsesTolerance(t *testing.T) {
	if Breached(10+FloatTolerance/2, 10, model.ComparatorAbove) {
		t.Fatal("values within tolerance must compare equal")
	}
}
