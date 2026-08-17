package store

import (
	"errors"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/zhangkui/go-alert-evaluator/internal/model"
)

var ErrInvalidMetric = errors.New("metric is required")
var ErrInvalidTimestamp = errors.New("timestamp is outside the accepted range")
var ErrInvalidValue = errors.New("sample value must be finite")

type SeriesStore struct {
	mu        sync.RWMutex
	retention time.Duration
	series    map[model.SeriesKey][]model.Sample
}

func NewSeriesStore(retention time.Duration) *SeriesStore {
	return &SeriesStore{retention: retention, series: make(map[model.SeriesKey][]model.Sample)}
}

func CanonicalLabels(labels model.Labels) string {
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(labels[key])
		builder.WriteByte(0)
	}
	return builder.String()
}

func (s *SeriesStore) Write(metric string, labels model.Labels, sample model.Sample, now time.Time) error {
	if strings.TrimSpace(metric) == "" {
		return ErrInvalidMetric
	}
	if math.IsNaN(sample.Value) || math.IsInf(sample.Value, 0) {
		return ErrInvalidValue
	}
	if sample.Timestamp.After(now.Add(time.Minute)) || sample.Timestamp.Before(now.Add(-s.retention)) {
		return ErrInvalidTimestamp
	}
	key := model.SeriesKey{Metric: metric, Labels: CanonicalLabels(labels)}
	s.mu.Lock()
	defer s.mu.Unlock()
	points := append(s.series[key], sample)
	sort.Slice(points, func(i, j int) bool { return points[i].Timestamp.Before(points[j].Timestamp) })
	s.series[key] = points
	return nil
}

func (s *SeriesStore) Window(metric string, labels model.Labels, start, end time.Time) []model.Sample {
	key := model.SeriesKey{Metric: metric, Labels: CanonicalLabels(labels)}
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Sample, 0, len(s.series[key]))
	for _, point := range s.series[key] {
		if !point.Timestamp.Before(start) && !point.Timestamp.After(end) {
			result = append(result, point)
		}
	}
	return result
}

func (s *SeriesStore) Cleanup(now time.Time) int {
	cutoff := now.Add(-s.retention)
	s.mu.Lock()
	defer s.mu.Unlock()
	removed := 0
	for key, points := range s.series {
		first := sort.Search(len(points), func(i int) bool { return !points[i].Timestamp.Before(cutoff) })
		removed += first
		if first == len(points) {
			delete(s.series, key)
			continue
		}
		s.series[key] = append([]model.Sample(nil), points[first:]...)
	}
	return removed
}
