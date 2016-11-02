// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// timeCmd times an external command.
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "time an external command",
	Long: `time runs an external command and reports the elapsed time in milliseconds 
to StatHat.`,
	Example: `  stathat time --name "rsync backup time" rsync -av ~/docs /opt/docs`,
	RunE:    runTime,
}

var statName string
var command string
var countErrs bool
var execCount int
var errSuffix string
var showOutput bool

func init() {
	RootCmd.AddCommand(timeCmd)
	timeCmd.Flags().StringVar(&statName, "name", "", "name of stat to report execution time")
	timeCmd.Flags().StringVar(&command, "exec", "", "the command to execute")
	timeCmd.Flags().BoolVar(&countErrs, "count-errs", true, "report error count to StatHat")
	timeCmd.Flags().StringVar(&errSuffix, "err-suffix", " - errors", "suffix added to `name` for error count stat")
	timeCmd.Flags().IntVar(&execCount, "n", 1, "number of times to run command")
	timeCmd.Flags().BoolVar(&showOutput, "output", false, "show command output")
	timeCmd.MarkFlagRequired("name")
}

func runTime(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && len(command) == 0 {
		return cmd.Usage()
	}
	if len(statName) == 0 {
		return cmd.Usage()
	}

	var allOut []byte
	if len(args) > 0 {
		command = strings.Join(args, " ")
	} else {
		args = strings.Fields(command)
	}
	start := time.Now()
	for i := 0; i < execCount; i++ {
		extCmd := exec.Command(args[0], args[1:]...)
		output, err := extCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("error running %q: %s", command, err)
			if countErrs {
				arg := ezarg{statName: statName + errSuffix, parameter: "count", value: 1}
				if perr := arg.post(); perr != nil {
					fmt.Printf("error posting error count: %s\n", perr)
				}
			}
			return err
		}
		allOut = append(allOut, output...)
	}

	elapsed := time.Since(start)
	ms := float64(elapsed) / float64(time.Millisecond)
	fmt.Printf("%s => %f ms (%s)\n", command, ms, elapsed)
	if showOutput {
		fmt.Printf("output:\n%s\n", allOut)
	}

	arg := ezarg{statName: statName, parameter: "value", value: ms}
	return arg.post()
}
