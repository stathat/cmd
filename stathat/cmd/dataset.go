// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"time"

	"github.com/stathat/cmd/stathat/db"
	"github.com/stathat/cmd/stathat/intr"

	"github.com/spf13/cobra"
)

// datasetCmd gets a dataset for a stat.
var datasetCmd = &cobra.Command{
	Use:   "dataset",
	Short: "get datasets for stats",
	RunE:  runDataset,
}

var timeframe string

func init() {
	RootCmd.AddCommand(datasetCmd)
	datasetCmd.Flags().BoolVar(&listJSON, "json", false, "display output as JSON")
	datasetCmd.Flags().BoolVar(&listCSV, "csv", false, "display output as CSV")
	datasetCmd.Flags().StringVar(&timeframe, "tf", "1w3h", "timeframe")
}

func runDataset(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || len(args) > 5 {
		return cmd.Usage()
	}

	store, err := db.New()
	if err != nil {
		return err
	}

	var datasets []intr.Dataset

	for _, id := range args {
		stat, ok := store.Lookup(id)
		if !ok {
			return fmt.Errorf("no stat found for %q", id)
		}
		dset, err := intr.LoadDataset(stat.ID, timeframe)
		if err != nil {
			return err
		}
		datasets = append(datasets, dset)
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

	return Output(&DataTable{dsets: datasets}, enc)
}

type DataTable struct {
	dsets []intr.Dataset
}

func (s *DataTable) Columns(row int) []string {
	t := time.Unix(s.dsets[0].Points[row].Time, 0)
	r := []string{t.Format(time.Stamp)}
	for _, d := range s.dsets {
		r = append(r, ftoa(d.Points[row].Value))
	}
	return r
}

func (s *DataTable) Header() []string {
	h := []string{"Time"}
	for _, d := range s.dsets {
		h = append(h, d.Name)
	}
	return h
}

func (s *DataTable) Len() int {
	if len(s.dsets) == 0 {
		return 0
	}
	return len(s.dsets[0].Points)
}

func (s *DataTable) Raw() interface{} {
	return s.dsets
}
