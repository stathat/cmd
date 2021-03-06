// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/stathat/cmd/stathat/intr"

	"github.com/spf13/cobra"
)

var listJSON, listCSV bool

// listCmd lists all the stats in an account.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "display all the stats in your account",
	RunE:  runList,
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&listJSON, "json", false, "display output as JSON")
	listCmd.Flags().BoolVar(&listCSV, "csv", false, "display output as CSV")
}

func runList(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		cmd.Usage()
		return nil
	}
	if listJSON && listCSV {
		cmd.Usage()
		return nil
	}

	stats, err := intr.StatList()
	if err != nil {
		return err
	}

	return outputStats(stats)
}

// StatTable represents a group of stats.  It contains
// methods to help output the data in a tabular format.
type StatTable struct {
	stats []intr.Stat
}

// NewStatTable creates a StatTable from a list of stats.
func NewStatTable(stats []intr.Stat) *StatTable {
	return &StatTable{stats: stats}
}

// Columns returns the string representation of the
// columns in a row.
func (s *StatTable) Columns(row int) []string {
	return s.stats[row].Strings()
}

// Header returns the column headers for the table.
func (s *StatTable) Header() []string {
	return []string{"ID", "Name", "Type", "Access"}
}

// Len returns the number of stats in this StatTable.
func (s *StatTable) Len() int {
	return len(s.stats)
}

// Raw returns a raw representation of the StatTable
// data.
func (s *StatTable) Raw() interface{} {
	return s.stats
}

func outputStats(stats []intr.Stat) error {
	t := NewStatTable(stats)
	var enc OutputEncoding
	switch {
	case listJSON:
		enc = OutputJSON
	case listCSV:
		enc = OutputCSV
	default:
		enc = OutputTab
	}

	return Output(t, enc)
}
