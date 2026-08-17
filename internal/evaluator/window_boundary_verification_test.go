package evaluator

import (
	"testing"
	"time"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
)

func TestAggregateIncludesSampleAtWindowStart(t *testing.T) {
	start := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	samples := []model.Sample{{Timestamp: start, Value: 2}, {Timestamp: start.Add(time.Minute), Value: 8}}
	result, err := Aggregate(model.AggregationAverage, samples, start)
	if err != nil { t.Fatal(err) }
	if result.Value != 5 { t.Fatalf("average = %v, want 5 with inclusive window start", result.Value) }
}
