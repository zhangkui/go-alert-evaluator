package store

import (
	"math"
	"testing"
	"time"

	"github.com/zhangkui/go-alert-evaluator/internal/model"
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

func TestSeriesStoreReplacesOutOfOrderSampleWithCanonicalLabels(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	series := NewSeriesStore(time.Hour)
	timestamps := []time.Time{
		now.Add(-3 * time.Minute),
		now.Add(-time.Minute),
		now.Add(-2 * time.Minute),
	}
	for index, timestamp := range timestamps {
		if err := series.Write("cpu", model.Labels{"host": "api-1", "zone": "east"}, model.Sample{Timestamp: timestamp, Value: float64(index + 1)}, now); err != nil {
			t.Fatal(err)
		}
	}

	if err := series.Write("cpu", model.Labels{"zone": "east", "host": "api-1"}, model.Sample{Timestamp: timestamps[2], Value: 34}, now); err != nil {
		t.Fatal(err)
	}

	points := series.Window("cpu", model.Labels{"host": "api-1", "zone": "east"}, now.Add(-5*time.Minute), now)
	if len(points) != 3 {
		t.Fatalf("got %d points, want 3: %#v", len(points), points)
	}
	for index := 1; index < len(points); index++ {
		if !points[index-1].Timestamp.Before(points[index].Timestamp) {
			t.Fatalf("samples are not strictly ordered: %#v", points)
		}
	}
	if points[1].Value != 34 {
		t.Fatalf("replacement value = %v, want 34", points[1].Value)
	}
}

func TestSeriesStoreRejectsInvalidSamplesWithoutReplacing(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	series := NewSeriesStore(time.Hour)
	timestamp := now.Add(-time.Minute)
	if err := series.Write("cpu", nil, model.Sample{Timestamp: timestamp, Value: 12}, now); err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name   string
		sample model.Sample
		want   error
	}{
		{name: "expired timestamp", sample: model.Sample{Timestamp: now.Add(-time.Hour - time.Nanosecond), Value: 34}, want: ErrInvalidTimestamp},
		{name: "NaN", sample: model.Sample{Timestamp: timestamp, Value: math.NaN()}, want: ErrInvalidValue},
		{name: "positive infinity", sample: model.Sample{Timestamp: timestamp, Value: math.Inf(1)}, want: ErrInvalidValue},
		{name: "negative infinity", sample: model.Sample{Timestamp: timestamp, Value: math.Inf(-1)}, want: ErrInvalidValue},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := series.Write("cpu", nil, test.sample, now); err != test.want {
				t.Fatalf("got error %v, want %v", err, test.want)
			}
		})
	}

	points := series.Window("cpu", nil, now.Add(-5*time.Minute), now)
	if len(points) != 1 || points[0].Value != 12 {
		t.Fatalf("invalid write changed stored points: %#v", points)
	}
}
