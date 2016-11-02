package intr

import "github.com/stathat/cmd/stathat/net"

type Stat struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Public  bool   `json:"public"`
	Counter bool   `json:"counter"`
}

func (s Stat) Kind() string {
	if s.Counter {
		return "counter"
	}
	return "value"
}

func (s Stat) Access() string {
	if s.Public {
		return "public"
	}
	return "private"
}

func (s Stat) Strings() []string {
	return []string{s.ID, s.Name, s.Kind(), s.Access()}
}

func StatList() ([]Stat, error) {
	var stats []Stat
	if err := net.DefaultAPI.Get("stats", nil, &stats); err != nil {
		return nil, err
	}
	return stats, nil
}
