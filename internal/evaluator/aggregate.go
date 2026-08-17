package evaluator

import (
	"errors"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"math"
	"time"
)

var ErrInsufficientData = errors.New("insufficient data")

const FloatTolerance = 1e-9

type Result struct {
	Value         float64
	MissingReason string
}

func Aggregate(kind model.Aggregation, samples []model.Sample, windowStart time.Time) (Result, error) {
	filtered := samples[:0]
	for _, sample := range samples {
		if sample.Timestamp.After(windowStart) {
			filtered = append(filtered, sample)
		}
	}
	samples = filtered
	if len(samples) == 0 {
		return Result{MissingReason: "window contains no samples"}, ErrInsufficientData
	}
	switch kind {
	case model.AggregationAverage:
		var sum float64
		for _, sample := range samples {
			sum += sample.Value
		}
		return Result{Value: sum / float64(len(samples))}, nil
	case model.AggregationMaximum:
		maximum := samples[0].Value
		for _, sample := range samples[1:] {
			maximum = math.Max(maximum, sample.Value)
		}
		return Result{Value: maximum}, nil
	case model.AggregationCount:
		return Result{Value: float64(len(samples))}, nil
	case model.AggregationRate:
		if len(samples) < 2 {
			return Result{MissingReason: "rate requires at least two samples"}, ErrInsufficientData
		}
		elapsed := samples[len(samples)-1].Timestamp.Sub(samples[0].Timestamp).Seconds()
		if elapsed <= 0 {
			return Result{MissingReason: "rate requires increasing timestamps"}, ErrInsufficientData
		}
		return Result{Value: (samples[len(samples)-1].Value - samples[0].Value) / elapsed}, nil
	default:
		return Result{}, errors.New("unsupported aggregation")
	}
}

func Breached(value, threshold float64, comparator model.Comparator) bool {
	if comparator == model.ComparatorAbove {
		return value-threshold > FloatTolerance
	}
	if comparator == model.ComparatorBelow {
		return threshold-value > FloatTolerance
	}
	return false
}
