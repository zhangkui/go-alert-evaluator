package store

import (
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"testing"
	"time"
)

func TestSeriesStoreOrdersSamples(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	series := NewSeriesStore(time.Hour)
	labels := model.Labels{"host": "api-1"}
	for _, offset := range []time.Duration{-time.Minute, -3 * time.Minute, -2 * time.Minute} {
		if err := series.Write("cpu", labels, model.Sample{Timestamp: now.Add(offset), Value: 1}, now); err != nil {
			t.Fatal(err)
		}
	}
	points := series.Window("cpu", labels, now.Add(-5*time.Minute), now)
	if len(points) != 3 || !points[0].Timestamp.Before(points[1].Timestamp) || !points[1].Timestamp.Before(points[2].Timestamp) {
		t.Fatalf("samples are not ordered: %#v", points)
	}
}

func TestSeriesStoreRejectsFutureSamples(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	err := NewSeriesStore(time.Hour).Write("cpu", nil, model.Sample{Timestamp: now.Add(2 * time.Minute), Value: 1}, now)
	if err != ErrInvalidTimestamp {
		t.Fatalf("expected invalid timestamp, got %v", err)
	}
}
