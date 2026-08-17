package silence

import (
	"errors"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"sync"
	"time"
)

var ErrInvalidInterval = errors.New("silence end must be after start")

type Store struct {
	mu       sync.RWMutex
	silences map[string]model.Silence
}

func NewStore() *Store { return &Store{silences: make(map[string]model.Silence)} }

func (s *Store) Add(item model.Silence) error {
	if item.ID == "" || !item.EndsAt.After(item.StartsAt) {
		return ErrInvalidInterval
	}
	item.Matchers = item.Matchers.Clone()
	s.mu.Lock()
	s.silences[item.ID] = item
	s.mu.Unlock()
	return nil
}

func (s *Store) Active(labels model.Labels, at time.Time) (model.Silence, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.silences {
		if at.Before(item.StartsAt) || !at.Before(item.EndsAt) {
			continue
		}
		matched := true
		for key, expected := range item.Matchers {
			if labels[key] != expected {
				matched = false
				break
			}
		}
		if matched {
			return item, true
		}
	}
	return model.Silence{}, false
}
