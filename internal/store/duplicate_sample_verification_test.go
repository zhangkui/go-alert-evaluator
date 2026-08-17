package store

import (
	"testing"
	"time"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
)

func TestDuplicateTimestampReplacesExistingSample(t *testing.T) {
	now := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	series := NewSeriesStore(time.Hour)
	labels := model.Labels{"host": "api-1"}
	for _, value := range []float64{12, 34} {
		if err := series.Write("cpu_usage", labels, model.Sample{Timestamp: now.Add(-time.Minute), Value: value}, now); err != nil { t.Fatal(err) }
	}
	points := series.Window("cpu_usage", labels, now.Add(-5*time.Minute), now)
	if len(points) != 1 { t.Fatalf("duplicate timestamp produced %d points, want 1", len(points)) }
	if points[0].Value != 34 { t.Fatalf("stored value = %v, want latest value 34", points[0].Value) }
}
