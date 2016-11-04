// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stathat/cmd/stathat/db"
	"github.com/stathat/cmd/stathat/intr"
)

// commentCmd adds a comment to a stat.
var commentCmd = &cobra.Command{
	Use:   `comment [stat id or name] -m "comment text"`,
	Short: "add a comment to a stat",
	Long:  `Comment adds text comments to stats.  You can view them with 'info'`,
	RunE:  runComment,
}

var message string

func init() {
	RootCmd.AddCommand(commentCmd)
	commentCmd.Flags().StringVar(&message, "m", "", "comment text")
}

func runComment(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Usage()
	}
	if len(message) == 0 {
		return errors.New("empty comment text")
	}

	store, err := db.New()
	if err != nil {
		return err
	}
	stat, ok := store.Lookup(args[0])
	if !ok {
		return fmt.Errorf("no stat found for %q (%d stats)", args[0], store.Count())
	}
	res, err := intr.AddComment(stat.ID, message)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
