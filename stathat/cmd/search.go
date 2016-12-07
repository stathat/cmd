package cmd

import (
	"fmt"
	"os"
	"strings"

	isatty "github.com/mattn/go-isatty"
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
	searchCmd.Flags().BoolVar(&listJSON, "json", false, "display output as JSON")
	searchCmd.Flags().BoolVar(&listCSV, "csv", false, "display output as CSV")
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

	if isatty.IsTerminal(os.Stdout.Fd()) || listJSON || listCSV {
		return outputStats(match)
	}

	ids := make([]string, len(match))
	for i, s := range match {
		ids[i] = s.ID
	}
	fmt.Println(strings.Join(ids, " "))

	return nil
}
