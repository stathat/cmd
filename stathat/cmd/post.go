// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "post stat data to StatHat",
}

var postValueCmd = &cobra.Command{
	Use:   "value",
	Short: "post value data point to StatHat",
	RunE:  runPostValue,
}

var postCountCmd = &cobra.Command{
	Use:   "count",
	Short: "post count data point to StatHat",
	RunE:  runPostCount,
}

func init() {
	RootCmd.AddCommand(postCmd)
	postCmd.AddCommand(postValueCmd)
	postCmd.AddCommand(postCountCmd)
}

func runPostValue(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		cmd.Usage()
	}

	value, err := atof(args[1])
	if err != nil {
		return fmt.Errorf("invalid value argument: %s", err)
	}

	arg := ezarg{statName: args[0], parameter: "value", value: value}
	return arg.post()
}

func runPostCount(cmd *cobra.Command, args []string) error {
	if len(args) != 1 && len(args) != 2 {
		cmd.Usage()
	}

	count := 1.0
	if len(args) == 2 {
		c, err := atof(args[1])
		if err != nil {
			return fmt.Errorf("invalid count argument: %s", err)
		}
		count = c
	}

	arg := ezarg{statName: args[0], parameter: "count", value: count}
	return arg.post()
}

type ezarg struct {
	statName  string
	parameter string
	value     float64
}

func (e ezarg) post() error {
	apiHost := viper.GetString("posthost")
	resp, err := http.PostForm(apiHost+"/ez",
		url.Values{
			"stat":      {e.statName},
			"ezkey":     {viper.GetString("ezkey")},
			e.parameter: {ftoa(e.value)},
		})
	if err != nil {
		return fmt.Errorf("post to %s error: %s", apiHost, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("post to %s failed: %s", apiHost, resp.Status)
	}

	fmt.Printf("worked %+v\n", resp)

	return nil
}
