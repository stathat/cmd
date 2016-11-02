// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "github.com/spf13/cobra"

// commentCmd adds a comment to a stat.
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "add a comment to a stat",
	RunE:  runComment,
}

func init() {
	RootCmd.AddCommand(commentCmd)
}

func runComment(cmd *cobra.Command, args []string) error {
	return nil
}
