package cmd

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/stathat/cmd/stathat/intr"
)

var correlationCmd = &cobra.Command{
	Use:   "correlation",
	Short: "calculate correlation matrix",
	RunE:  runCorrelation,
}

func init() {
	RootCmd.AddCommand(correlationCmd)
}

type crow struct {
	dset   intr.Dataset
	stat   intr.Stat
	mean   float64
	deltas []float64
	vari   float64
	corrs  map[string]float64
}

func runCorrelation(cmd *cobra.Command, args []string) error {
	stats, err := intr.StatList()
	if err != nil {
		return err
	}

	table := make(map[string]*crow)

	/*
		fmt.Println("limit of 100 stats during devel")
		if len(stats) > 100 {
			stats = stats[0:100]
		}
	*/

	for _, stat := range stats {
		fmt.Printf("loading dataset for %q\n", stat.Name)
		dset, err := intr.LoadDataset(stat.ID, "1w1h")
		if err != nil {
			return err
		}
		table[stat.ID] = &crow{
			dset:  dset,
			stat:  stat,
			corrs: make(map[string]float64),
		}
	}

	start := time.Now()

	for _, row := range table {
		num := len(row.dset.Points)
		if num == 0 {
			continue
		}
		var sum float64
		for _, p := range row.dset.Points {
			sum += p.Value
		}
		row.mean = sum / float64(num)

		row.deltas = make([]float64, num)
		sum = 0.0
		for i, p := range row.dset.Points {
			row.deltas[i] = p.Value - row.mean
			sum += (row.deltas[i] * row.deltas[i])
		}
		row.vari = math.Sqrt(sum)
	}

	nstats := len(stats)
	for i := 0; i < nstats; i++ {
		statI := stats[i]
		rowI := table[statI.ID]
		num := len(rowI.dset.Points)
		for j := i + 1; j < nstats; j++ {
			statJ := stats[j]
			rowJ := table[statJ.ID]

			if rowI.vari == 0.0 || rowJ.vari == 0.0 {
				rowI.corrs[statJ.ID] = 0.0
				rowJ.corrs[statI.ID] = 0.0
				continue
			}

			var total float64
			for k := 0; k < num; k++ {
				total += rowI.deltas[k] * rowJ.deltas[k]
			}

			corr := total / (rowI.vari * rowJ.vari)
			rowI.corrs[statJ.ID] = corr
			rowJ.corrs[statI.ID] = corr
		}
	}

	for k, row := range table {
		if len(row.corrs) == 0 {
			continue
		}
		sorted := sortMapByValue(row.corrs)
		if sorted[0].Value == 0.0 {
			continue
		}
		if len(sorted) > 10 {
			sorted = sorted[:10]
		}
		fmt.Printf("%s\t%.3f\t%.3f\t%s\n", k, row.mean, row.vari, row.stat.Name)
		for _, pair := range sorted {
			if pair.Value != 0.0 {
				fmt.Printf("\t%s => %.3f\n", pair.Key, pair.Value)
			}
		}
		fmt.Println("----------------------------------------------------------------")
	}

	elapsed := time.Since(start)
	fmt.Printf("correlation matrix time: %s", elapsed)

	return nil
}

type Pair struct {
	Key   string
	Value float64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return math.Abs(p[i].Value) > math.Abs(p[j].Value) }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]float64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}
