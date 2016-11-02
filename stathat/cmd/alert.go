// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "github.com/spf13/cobra"

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Manage alerts",
}

var alertValueCmd = &cobra.Command{
	Use:   "value",
	Short: "Create a value alert",
	RunE:  runAlertValue,
}

var alertDeltaCmd = &cobra.Command{
	Use:   "delta",
	Short: "Create a value alert",
	RunE:  runAlertDelta,
}

var alertDataCmd = &cobra.Command{
	Use:   "data",
	Short: "Create a value alert",
	RunE:  runAlertData,
}

var alertDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an alert",
	RunE:  runAlertDelete,
}

var alertListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all alerts",
	RunE:  runAlertList,
}

var alertInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about an alert",
	RunE:  runAlertInfo,
}

func init() {
	RootCmd.AddCommand(alertCmd)
	alertCmd.AddCommand(alertValueCmd)
	alertCmd.AddCommand(alertDeltaCmd)
	alertCmd.AddCommand(alertDataCmd)
	alertCmd.AddCommand(alertDeleteCmd)
	alertCmd.AddCommand(alertListCmd)
	alertCmd.AddCommand(alertInfoCmd)
}

func runAlertValue(cmd *cobra.Command, args []string) error {
	return nil
}
func runAlertDelta(cmd *cobra.Command, args []string) error {
	return nil
}
func runAlertData(cmd *cobra.Command, args []string) error {
	return nil
}
func runAlertDelete(cmd *cobra.Command, args []string) error {
	return nil
}
func runAlertList(cmd *cobra.Command, args []string) error {
	return nil
}
func runAlertInfo(cmd *cobra.Command, args []string) error {
	return nil
}
