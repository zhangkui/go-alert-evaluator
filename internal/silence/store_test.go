package silence

import (
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"testing"
	"time"
)

func TestActiveSilenceMatchesSingleLabel(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	store := NewStore()
	if err := store.Add(model.Silence{ID: "s1", Matchers: model.Labels{"service": "api"}, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}); err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Active(model.Labels{"service": "api", "region": "east"}, now); !ok {
		t.Fatal("expected active silence")
	}
}
