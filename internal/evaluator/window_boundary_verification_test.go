package evaluator

import (
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"math"
	"testing"
	"time"
)

func TestAggregateIncludesSampleAtWindowStart(t *testing.T) {
	start := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	samples := []model.Sample{
		{Timestamp: start.Add(-time.Nanosecond), Value: 100},
		{Timestamp: start, Value: 2},
		{Timestamp: start.Add(time.Minute), Value: 8},
	}
	tests := []struct {
		name string
		kind model.Aggregation
		want float64
	}{
		{name: "average", kind: model.AggregationAverage, want: 5},
		{name: "count", kind: model.AggregationCount, want: 2},
		{name: "maximum", kind: model.AggregationMaximum, want: 8},
		{name: "rate", kind: model.AggregationRate, want: 0.1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Aggregate(test.kind, append([]model.Sample(nil), samples...), start)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(result.Value-test.want) > FloatTolerance {
				t.Fatalf("%s = %v, want %v with inclusive window start", test.name, result.Value, test.want)
			}
		})
	}
}
