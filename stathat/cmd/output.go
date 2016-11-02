// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func atof(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}

type Table interface {
	Header() []string
	Columns(row int) []string
	Len() int
	Raw() interface{}
}

type OutputEncoding int

const (
	OutputRaw OutputEncoding = iota
	OutputJSON
	OutputCSV
	OutputTab
)

func Output(t Table, enc OutputEncoding) error {
	switch enc {
	case OutputRaw:
		return outputRaw(t)
	case OutputJSON:
		return outputJSON(t)
	case OutputCSV:
		return outputCSV(t)
	case OutputTab:
		return outputTab(t)
	default:
		return fmt.Errorf("invalid output encoding: %d", enc)
	}
}

func outputRaw(t Table) error {
	o, err := json.Marshal(t.Raw())
	if err != nil {
		return err
	}
	fmt.Println(string(o))
	return nil
}

func outputJSON(t Table) error {
	o, err := json.MarshalIndent(t.Raw(), "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(o))
	return nil
}

func outputCSV(t Table) error {
	w := csv.NewWriter(os.Stdout)
	w.Write(t.Header())
	for i := 0; i < t.Len(); i++ {
		w.Write(t.Columns(i))
	}
	w.Flush()
	return nil
}

func outputTab(t Table) error {
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, strings.Join(t.Header(), "\t"))
	for i := 0; i < t.Len(); i++ {
		fmt.Fprintf(w, strings.Join(t.Columns(i), "\t")+"\n")
	}
	w.Flush()
	return nil
}
