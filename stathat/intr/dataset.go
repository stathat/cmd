package intr

import (
	"net/url"
	"path"
	"strconv"

	"github.com/stathat/cmd/stathat/net"
)

type Dataset struct {
	Name      string
	Timeframe string
	Kind      string
	Points    []Point
}

type Point struct {
	Time   int64
	Value  float64
	Rcount float64
}

// XXX repeat
func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

// XXX repeat
func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}

func (p Point) Strings() []string {
	return []string{itoa(p.Time), ftoa(p.Value), ftoa(p.Rcount)}
}

func LoadDataset(id, timeframe string) (Dataset, error) {
	var dset Dataset
	p := path.Join("stats", id, "dataset")
	v := url.Values{}
	v.Set("t", timeframe)
	if err := net.DefaultAPI.Get(p, v, &dset); err != nil {
		return Dataset{}, err
	}
	return dset, nil
}

func LoadDatasetFull(id string) (Dataset, error) {
	var dset Dataset
	p := path.Join("stats", id, "dataset_full")
	v := url.Values{}
	if err := net.DefaultAPI.Get(p, v, &dset); err != nil {
		return Dataset{}, err
	}
	return dset, nil
}
