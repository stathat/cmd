package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/stathat/cmd/stathat/intr"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "search for stats",
	RunE:  runSearch,
}

func init() {
	RootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		cmd.Usage()
		return nil
	}
	query := strings.ToLower(args[0])

	stats, err := intr.StatList()
	if err != nil {
		return err
	}

	var match []intr.Stat
	for _, s := range stats {
		if !strings.Contains(strings.ToLower(s.Name), query) {
			continue
		}
		match = append(match, s)
	}

	return outputStats(match)
}
