package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/spf13/cobra"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "display and update sorted information about stats",
	RunE:  top,
}

func init() {
	RootCmd.AddCommand(topCmd)
}

type mainWindow struct {
	parent *views.Application
	view   views.View
	main   *views.CellView
	keybar *views.SimpleStyledText
	status *views.SimpleStyledTextBar

	views.Panel
}

func (a *mainWindow) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'Q', 'q':
				a.parent.Quit()
				return true
			}
		}
	}

	return a.Panel.HandleEvent(ev)
}

func clock(app *views.Application, title *views.TextBar) {
	for {
		title.SetRight(time.Now().Format("15:04:05"), tcell.StyleDefault)
		app.Refresh()
		time.Sleep(1 * time.Second)
	}
}

func top(cmd *cobra.Command, args []string) error {
	app := &views.Application{}
	window := &mainWindow{parent: app}

	title := views.NewTextBar()
	title.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorTeal).
		Foreground(tcell.ColorWhite))
	title.SetCenter("StatHat top", tcell.StyleDefault)
	title.SetRight("21:22:37", tcell.StyleDefault)

	go clock(app, title)

	window.SetTitle(title)

	app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	app.SetRootWidget(window)

	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	return nil
}
