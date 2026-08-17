package engine

import (
	"errors"
	"github.com/zhangkui/go-alert-evaluator/internal/evaluator"
	"github.com/zhangkui/go-alert-evaluator/internal/model"
	"github.com/zhangkui/go-alert-evaluator/internal/silence"
	"github.com/zhangkui/go-alert-evaluator/internal/store"
	"sync"
	"time"
)

var ErrRuleNotFound = errors.New("rule not found")

type Sender interface {
	Send(model.Rule, model.Evaluation) error
}
type NopSender struct{}

func (NopSender) Send(model.Rule, model.Evaluation) error { return nil }

type runtimeState struct {
	state        model.AlertState
	pendingSince time.Time
	recoverSince time.Time
}

type Engine struct {
	mu       sync.Mutex
	series   *store.SeriesStore
	silences *silence.Store
	sender   Sender
	rules    map[string]model.Rule
	states   map[string]runtimeState
	history  map[string][]model.Evaluation
}

func New(series *store.SeriesStore, silences *silence.Store, sender Sender) *Engine {
	if sender == nil {
		sender = NopSender{}
	}
	return &Engine{series: series, silences: silences, sender: sender, rules: make(map[string]model.Rule), states: make(map[string]runtimeState), history: make(map[string][]model.Evaluation)}
}

func (e *Engine) AddRule(rule model.Rule) error {
	if rule.ID == "" || rule.Metric == "" || rule.Window <= 0 || rule.For < 0 || rule.RecoverFor < 0 {
		return errors.New("invalid rule")
	}
	rule.Labels = rule.Labels.Clone()
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules[rule.ID] = rule
	if _, ok := e.states[rule.ID]; !ok {
		e.states[rule.ID] = runtimeState{state: model.StateNormal}
	}
	return nil
}

func (e *Engine) Evaluate(ruleID string, at time.Time) (model.Evaluation, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	rule, ok := e.rules[ruleID]
	if !ok {
		return model.Evaluation{}, ErrRuleNotFound
	}
	current := e.states[ruleID]
	previous := current.state
	result, err := evaluator.Aggregate(rule.Aggregation, e.series.Window(rule.Metric, rule.Labels, at.Add(-rule.Window), at), at.Add(-rule.Window))
	evaluation := model.Evaluation{RuleID: ruleID, At: at, PreviousState: previous, Notification: "none"}
	if err != nil {
		current = runtimeState{state: model.StateInsufficient}
		evaluation.State = current.state
		evaluation.MissingReason = result.MissingReason
		e.finish(rule, current, &evaluation)
		return evaluation, nil
	}
	evaluation.Value = new(float64)
	*evaluation.Value = result.Value
	if evaluator.Breached(result.Value, rule.Threshold, rule.Comparator) {
		current.recoverSince = time.Time{}
		switch current.state {
		case model.StateNormal, model.StateInsufficient:
			if rule.For == 0 {
				current.state = model.StateFiring
			} else {
				current.state = model.StatePending
				current.pendingSince = at
			}
		case model.StatePending:
			if at.Sub(current.pendingSince) >= rule.For {
				current.state = model.StateFiring
			}
		}
	} else {
		current.pendingSince = time.Time{}
		if current.state == model.StateFiring && rule.RecoverFor > 0 {
			if current.recoverSince.IsZero() {
				current.recoverSince = at
			}
			current.state = model.StateNormal
		} else {
			current.state = model.StateNormal
			current.recoverSince = time.Time{}
		}
	}
	evaluation.State = current.state
	e.finish(rule, current, &evaluation)
	return evaluation, nil
}

func (e *Engine) finish(rule model.Rule, current runtimeState, evaluation *model.Evaluation) {
	if evaluation.PreviousState != evaluation.State {
		evaluation.Notification = string(evaluation.State)
		if _, muted := e.silences.Active(rule.Labels, evaluation.At); muted {
			evaluation.Notification = "suppressed"
		} else if err := e.sender.Send(rule, *evaluation); err == nil {
			evaluation.NotificationSent = true
		} else {
			evaluation.Notification = "send_failed"
		}
	}
	e.states[rule.ID] = current
	e.history[rule.ID] = append(e.history[rule.ID], *evaluation)
}

func (e *Engine) History(ruleID string) []model.Evaluation {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]model.Evaluation(nil), e.history[ruleID]...)
}
