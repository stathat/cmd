// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "github.com/spf13/cobra"

// deleteCmd deletes a stat.
var deleteCmd = &cobra.Command{
	Use:    "delete",
	Short:  "delete a stat",
	RunE:   runDelete,
	Hidden: true,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	return nil
}
