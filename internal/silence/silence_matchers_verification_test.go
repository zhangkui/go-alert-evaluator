package silence

import (
	"testing"
	"time"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
)

func TestSilenceRequiresEveryMatcher(t *testing.T) {
	now := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	store := NewStore()
	if err := store.Add(model.Silence{ID: "west-api", Matchers: model.Labels{"service": "api", "region": "west"}, StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Minute)}); err != nil { t.Fatal(err) }
	if _, ok := store.Active(model.Labels{"service": "api", "region": "east"}, now); ok { t.Fatal("silence matched even though region differed") }
}
