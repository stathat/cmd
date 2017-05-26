// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping StatHat",
	RunE:  runPing,
}

func init() {
	RootCmd.AddCommand(pingCmd)
}

func runPing(cmd *cobra.Command, args []string) error {
	apiHost := viper.GetString("posthost")
	resp, err := http.Get(apiHost)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ping to %s failed: %s", apiHost, resp.Status)
	}

	fmt.Printf("ping %s success.\n", apiHost)

	return nil
}
