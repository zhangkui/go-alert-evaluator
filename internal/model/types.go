package model

import "time"

type Labels map[string]string

func (l Labels) Clone() Labels {
	clone := make(Labels, len(l))
	for key, value := range l {
		clone[key] = value
	}
	return clone
}

type Sample struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}
type SeriesKey struct {
	Metric string
	Labels string
}
type Aggregation string

const (
	AggregationAverage Aggregation = "avg"
	AggregationMaximum Aggregation = "max"
	AggregationCount   Aggregation = "count"
	AggregationRate    Aggregation = "rate"
)

type Comparator string

const (
	ComparatorAbove Comparator = "above"
	ComparatorBelow Comparator = "below"
)

type Rule struct {
	ID          string        `json:"id"`
	Metric      string        `json:"metric"`
	Labels      Labels        `json:"labels"`
	Aggregation Aggregation   `json:"aggregation"`
	Comparator  Comparator    `json:"comparator"`
	Threshold   float64       `json:"threshold"`
	Window      time.Duration `json:"-"`
	For         time.Duration `json:"-"`
	RecoverFor  time.Duration `json:"-"`
}

type AlertState string

const (
	StateNormal       AlertState = "normal"
	StatePending      AlertState = "pending"
	StateFiring       AlertState = "firing"
	StateInsufficient AlertState = "insufficient_data"
)

type Evaluation struct {
	RuleID           string     `json:"rule_id"`
	At               time.Time  `json:"at"`
	PreviousState    AlertState `json:"previous_state"`
	State            AlertState `json:"state"`
	Value            *float64   `json:"value,omitempty"`
	MissingReason    string     `json:"missing_reason,omitempty"`
	Notification     string     `json:"notification"`
	NotificationSent bool       `json:"notification_sent"`
}

type Silence struct {
	ID       string    `json:"id"`
	Matchers Labels    `json:"matchers"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}
