// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"github.com/stathat/cmd/stathat/db"
	"github.com/stathat/cmd/stathat/intr"

	"github.com/spf13/cobra"
)

// summaryCmd displays data summary for a stat.
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "stat data summary",
	RunE:  runSummary,
}

func init() {
	RootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().BoolVar(&listJSON, "json", false, "display output as JSON")
	summaryCmd.Flags().BoolVar(&listCSV, "csv", false, "display output as JSON")
	summaryCmd.Flags().StringVar(&timeframe, "tf", "1w3h", "timeframe")
}

func runSummary(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || len(args) > 5 {
		return cmd.Usage()
	}

	store, err := db.New()
	if err != nil {
		return err
	}

	var summaries []*Summary
	for _, id := range args {
		stat, ok := store.Lookup(id)
		if !ok {
			return fmt.Errorf("no stat found for %q", id)
		}
		dset, err := intr.LoadDataset(stat.ID, timeframe)
		if err != nil {
			return err
		}
		summaries = append(summaries, NewSummary(&dset))
	}

	var enc OutputEncoding
	switch {
	case listJSON:
		enc = OutputJSON
	case listCSV:
		enc = OutputCSV
	default:
		enc = OutputTab
	}

	for _, s := range summaries {
		switch enc {
		case OutputRaw:
			o, err := json.Marshal(s)
			if err != nil {
				return err
			}
			fmt.Println(string(o))
		case OutputJSON:
			o, err := json.MarshalIndent(s, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println(string(o))
		case OutputTab:
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 2, '\t', 0)
			fmt.Fprintf(w, "Name\t%s\n", s.Name)
			fmt.Fprintf(w, "Latest\t%g\n", s.Latest)
			fmt.Fprintf(w, "Min\t%g\n", s.Min)
			fmt.Fprintf(w, "Max\t%g\n", s.Max)
			fmt.Fprintf(w, "Mean\t%g\n", s.Mean)
			if s.Counter {
				fmt.Fprintf(w, "Total\t%g\n", s.Total)
			}
			fmt.Fprintf(w, "Std Dev\t%g\n", s.StdDev)
			fmt.Fprintf(w, "95%%\t%g - %g\n", s.Conf95Min, s.Conf95Max)
			fmt.Fprintf(w, "99%%\t%g - %g\n", s.Conf99Min, s.Conf99Max)
			w.Flush()
		default:
			return fmt.Errorf("invalid output encoding: %d", enc)
		}
	}
	return nil
}

type Summary struct {
	Name      string
	Timeframe string
	Counter   bool
	Latest    float64
	Mean      float64
	Total     float64
	Min       float64
	Max       float64
	StdDev    float64
	Conf95Min float64
	Conf95Max float64
	Conf99Min float64
	Conf99Max float64
}

func NewSummary(d *intr.Dataset) *Summary {
	s := &Summary{Name: d.Name, Timeframe: d.Timeframe, Counter: d.Kind == "count"}
	var count int
	for _, p := range d.Points {
		if p.Rcount == 0 {
			continue
		}

		// point has data:

		if count == 0 {
			s.Min = p.Value
			s.Max = p.Value
		}

		if p.Value < s.Min {
			s.Min = p.Value
		}
		if p.Value > s.Max {
			s.Max = p.Value
		}

		s.Total += p.Value
		s.Latest = p.Value
		count++
	}
	if count == 0 {
		return s
	}

	s.Mean = s.Total / float64(count)

	if count == 1 {
		return s
	}

	var deltaSum float64
	for _, p := range d.Points {
		if p.Rcount == 0 {
			continue
		}
		delta := p.Value - s.Mean
		deltaSum += delta * delta
	}

	variance := deltaSum / float64(count-1)
	s.StdDev = math.Sqrt(variance)
	stderr := s.StdDev / math.Sqrt(float64(count))
	delta95 := 1.95 * stderr
	delta99 := 2.575 * stderr
	s.Conf95Min = s.Mean - delta95
	s.Conf95Max = s.Mean + delta95
	s.Conf99Min = s.Mean - delta99
	s.Conf99Max = s.Mean + delta99

	return s
}

type SummaryTable struct {
	summaries []*Summary
}

func (s *SummaryTable) Columns(row int) []string {
	return []string{"col 1", "col 2", "col 3"}
}

func (s *SummaryTable) Header() []string {
	return []string{"Name", "Min", "Max"}
}

func (s *SummaryTable) Len() int {
	return len(s.summaries)
}

func (s *SummaryTable) Raw() interface{} {
	return s.summaries
}
