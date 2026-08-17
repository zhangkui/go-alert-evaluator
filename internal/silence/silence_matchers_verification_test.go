package silence

import (
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"testing"
	"time"
)

func TestSilenceRequiresEveryMatcher(t *testing.T) {
	now := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	store := NewStore()
	if err := store.Add(model.Silence{ID: "west-api", Matchers: model.Labels{"service": "api", "region": "west"}, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}); err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Active(model.Labels{"service": "api", "region": "west"}, now); !ok {
		t.Fatal("expected silence to match all labels")
	}
	if _, ok := store.Active(model.Labels{"service": "api", "region": "east"}, now); ok {
		t.Fatal("silence matched even though region differed")
	}
	if _, ok := store.Active(model.Labels{"service": "web", "region": "west"}, now); ok {
		t.Fatal("silence matched even though service differed")
	}
}

func TestSilenceIntervalIsStartInclusiveEndExclusive(t *testing.T) {
	now := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	store := NewStore()
	if err := store.Add(model.Silence{ID: "s1", Matchers: model.Labels{"service": "api", "region": "west"}, StartsAt: now, EndsAt: now.Add(time.Minute)}); err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Active(model.Labels{"service": "api", "region": "west"}, now); !ok {
		t.Fatal("expected silence at start boundary")
	}
	if _, ok := store.Active(model.Labels{"service": "api", "region": "west"}, now.Add(time.Minute)); ok {
		t.Fatal("silence should be inactive at end boundary")
	}
}
