// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "set up StatHat access keys",
	Long:  `Setup configures the CLI access keys.`,
	RunE:  setup,
}

func setup(cmd *cobra.Command, args []string) error {
	var accessKey, ezKey string
	fmt.Printf("Access key:  ")
	n, err := fmt.Scanln(&accessKey)
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("stathat setup: invalid number of elts scanned")
	}
	fmt.Printf("EZ key:  ")
	fmt.Scanln(&ezKey)
	if err != nil {
		return fmt.Errorf("stathat setup: error getting ez key: %s", err)
	}
	if n != 1 {
		return errors.New("stathat setup: invalid number of elts scanned")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
