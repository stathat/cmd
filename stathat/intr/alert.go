package intr

import "github.com/stathat/cmd/stathat/net"

func AlertList() ([]Alert, error) {
	var alerts []Alert
	if err := net.DefaultAPI.Get("alerts", nil, &alerts); err != nil {
		return []Alert{}, err
	}
	return alerts, nil
}

type Alert struct {
	ID         string  `json:"id"`
	StatID     string  `json:"stat_id"`
	StatName   string  `json:"stat_name"`
	Kind       string  `json:"kind"`
	TimeWindow string  `json:"time_window"`
	Operator   string  `json:"operator,omitempty"`
	Threshold  float64 `json:"threshold,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
	TimeDelta  string  `json:"time_delta,omitempty"`
}

func (s Alert) Strings() []string {
	return []string{s.ID, s.StatID, s.StatName, s.Kind, s.TimeWindow}
}

type AlertTable struct {
	alerts []Alert
}

func NewAlertTable(alerts []Alert) *AlertTable {
	return &AlertTable{alerts: alerts}
}

func (s *AlertTable) Columns(row int) []string {
	return s.alerts[row].Strings()
}

func (s *AlertTable) Header() []string {
	return []string{"ID", "Stat ID", "Stat Name", "Kind", "Time Window"}
}

func (s *AlertTable) Len() int {
	return len(s.alerts)
}

func (s *AlertTable) Raw() interface{} {
	return s.alerts
}
