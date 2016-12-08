package cmd

import (
	"testing"
	"unicode"

	"github.com/spf13/cobra"
)

func check(t *testing.T, c *cobra.Command) {
	if len(c.Short) == 0 {
		t.Errorf("%s has no Short field", c.Name())
	}
	if unicode.IsUpper([]rune(c.Short)[0]) {
		t.Errorf("%s Short starts with uppercase: %q", c.Name(), c.Short)
	}
	if c.Short[len(c.Short)-1] == '.' {
		t.Errorf("%s Short ends with period: %s", c.Name(), c.Short)
	}
	for _, sub := range c.Commands() {
		check(t, sub)
	}
}

func TestShortFormatting(t *testing.T) {
	check(t, RootCmd)
}
