// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"net/url"

	"github.com/stathat/cmd/stathat/db"
	"github.com/stathat/cmd/stathat/intr"
	"github.com/stathat/cmd/stathat/net"

	"github.com/spf13/cobra"
)

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "manage alerts",
}

var alertValueCmd = &cobra.Command{
	Use:   "value",
	Short: "create a value alert",
	RunE:  runAlertValue,
}

var alertDeltaCmd = &cobra.Command{
	Use:   "delta",
	Short: "create a delta alert",
	RunE:  runAlertDelta,
}

var alertDataCmd = &cobra.Command{
	Use:   "data",
	Short: "create a data alert",
	RunE:  runAlertData,
}

var alertCreateFlags struct {
	TimeWindow string
	Percentage float64
	Operator   string
	TimeDelta  string
	Threshold  float64
}

var alertDeleteCmd = &cobra.Command{
	Use:    "delete",
	Short:  "delete an alert",
	RunE:   runAlertDelete,
	Hidden: true,
}

var alertListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all alerts",
	RunE:  runAlertList,
}

var alertInfoCmd = &cobra.Command{
	Use:    "info",
	Short:  "get info about an alert",
	RunE:   runAlertInfo,
	Hidden: true,
}

func init() {
	RootCmd.AddCommand(alertCmd)
	alertCmd.AddCommand(alertDeleteCmd)
	alertCmd.AddCommand(alertInfoCmd)

	alertCmd.AddCommand(alertListCmd)
	alertListCmd.Flags().BoolVar(&listJSON, "json", false, "display output as JSON")
	alertListCmd.Flags().BoolVar(&listCSV, "csv", false, "display output as CSV")

	alertCmd.AddCommand(alertDataCmd)
	addAlertCommonFlags(alertDataCmd)

	alertCmd.AddCommand(alertDeltaCmd)
	addAlertCommonFlags(alertDeltaCmd)
	alertDeltaCmd.Flags().Float64Var(&alertCreateFlags.Percentage, "percentage", 50.0, "change percentage")
	alertDeltaCmd.Flags().StringVar(&alertCreateFlags.Operator, "operator", "gt", "alert operator [gt, lt, dt]")
	alertDeltaCmd.Flags().StringVar(&alertCreateFlags.TimeDelta, "time_delta", "1d", "compare to time difference")

	alertCmd.AddCommand(alertValueCmd)
	addAlertCommonFlags(alertValueCmd)
}

func runAlertValue(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Usage()
	}
	store, err := db.New()
	if err != nil {
		return err
	}
	stat, ok := store.Lookup(args[0])
	if !ok {
		return fmt.Errorf("no stat found for %q", args[0])
	}

	p := "alerts"
	v := url.Values{}
	v.Set("kind", "value")
	v.Set("stat_id", stat.ID)
	v.Set("time_window", alertCreateFlags.TimeWindow)
	v.Set("threshold", ftoa(alertCreateFlags.Threshold))
	v.Set("operator", alertCreateFlags.Operator)
	var x interface{}
	err = net.DefaultAPI.Post(p, v, &x)
	return err
}

func runAlertDelta(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Usage()
	}
	store, err := db.New()
	if err != nil {
		return err
	}
	stat, ok := store.Lookup(args[0])
	if !ok {
		return fmt.Errorf("no stat found for %q", args[0])
	}

	p := "alerts"
	v := url.Values{}
	v.Set("kind", "delta")
	v.Set("stat_id", stat.ID)
	v.Set("time_window", alertCreateFlags.TimeWindow)
	v.Set("operator", alertCreateFlags.Operator)
	v.Set("time_delta", alertCreateFlags.TimeDelta)
	v.Set("percentage", ftoa(alertCreateFlags.Percentage))
	var x interface{}
	err = net.DefaultAPI.Post(p, v, &x)
	return err
}

func runAlertData(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Usage()
	}
	store, err := db.New()
	if err != nil {
		return err
	}
	stat, ok := store.Lookup(args[0])
	if !ok {
		return fmt.Errorf("no stat found for %q", args[0])
	}

	p := "alerts"
	v := url.Values{}
	v.Set("kind", "data")
	v.Set("stat_id", stat.ID)
	v.Set("time_window", alertCreateFlags.TimeWindow)
	var x interface{}
	err = net.DefaultAPI.Post(p, v, &x)
	return err
}

func runAlertDelete(cmd *cobra.Command, args []string) error {
	return nil
}

func runAlertList(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Usage()
	}

	alerts, err := intr.AlertList()
	if err != nil {
		return err
	}

	return outputAlerts(alerts)
}
func runAlertInfo(cmd *cobra.Command, args []string) error {
	return nil
}

func outputAlerts(alerts []intr.Alert) error {
	t := intr.NewAlertTable(alerts)

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

func addAlertCommonFlags(c *cobra.Command) {
	c.Flags().StringVar(&alertCreateFlags.TimeWindow, "time_window", "1d", "time window")
}
