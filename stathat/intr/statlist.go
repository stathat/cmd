package intr

import "github.com/stathat/cmd/stathat/net"

type Stat struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Public         bool   `json:"public"`
	Counter        bool   `json:"counter"`
	DataReceivedAt int    `json:"data_received_at"`
	CreatedAt      int    `json:"created_at"`
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

type ByDataReceivedAt []Stat

func (a ByDataReceivedAt) Len() int           { return len(a) }
func (a ByDataReceivedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDataReceivedAt) Less(i, j int) bool { return a[i].DataReceivedAt > a[j].DataReceivedAt }

type ByCreatedAt []Stat

func (a ByCreatedAt) Len() int           { return len(a) }
func (a ByCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedAt) Less(i, j int) bool { return a[i].CreatedAt > a[j].CreatedAt }
