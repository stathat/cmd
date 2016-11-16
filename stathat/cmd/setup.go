// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "set up StatHat access keys",
	Long:  `Setup configures the CLI access keys.`,
	RunE:  setup,
}

var setupForce bool

func init() {
	RootCmd.AddCommand(setupCmd)
	setupCmd.Flags().BoolVar(&setupForce, "force", false, "force setup (overwrite existing config)")
}

func setup(cmd *cobra.Command, args []string) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	filename := filepath.Join(u.HomeDir, ".stathat", "config.yaml")
	if !setupForce {
		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			return fmt.Errorf("config file %s already exists (use --force to overwrite)", filename)
		}
	}

	var accessKey, ezKey string
	fmt.Printf("Please enter an Access Token for your StatHat account.\n")
	fmt.Printf("You can get one here: https://www.stathat.com/access\n\n")
	fmt.Printf("Access token:  ")
	n, err := fmt.Scanln(&accessKey)
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("stathat setup: invalid number of elts scanned")
	}
	fmt.Printf("Please enter the EZ Key for your StatHat account.\n")
	fmt.Printf("You can find it (and change it) on the settings page:\n\n")
	fmt.Printf("\thttps://www.stathat.com/settings\n\n")
	fmt.Printf("EZ Key:  ")
	fmt.Scanln(&ezKey)
	if err != nil {
		return fmt.Errorf("stathat setup: error getting ez key: %s", err)
	}
	if n != 1 {
		return errors.New("stathat setup: invalid number of elts scanned")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "accesskey: %s\nezkey: %s\n", accessKey, ezKey)
	f.Close()

	return nil
}
