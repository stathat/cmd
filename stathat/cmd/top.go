package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/spf13/cobra"
	"github.com/stathat/cmd/stathat/config"
	"github.com/stathat/cmd/stathat/intr"
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

var (
	StyleNormal = tcell.StyleDefault.Foreground(tcell.ColorSilver).Background(tcell.ColorBlack)
	StyleGood   = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	StyleWarn   = tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack)
	StyleError  = tcell.StyleDefault.Foreground(tcell.ColorMaroon).Background(tcell.ColorBlack)
)

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

type model struct {
	stats     []intr.Stat
	summaries map[string]string
	sync.Mutex
}

func newModel() *model {
	return &model{
		summaries: make(map[string]string),
	}
}

func (m *model) SetStats(s []intr.Stat) {
	m.Lock()
	fmt.Printf("set stats: %v", s)
	m.stats = s
	m.Unlock()
}

func (m *model) SetSummary(id, summary string) {
	m.Lock()
	m.summaries[id] = summary
	m.Unlock()
}

func (m *model) GetBounds() (int, int) {
	return 80, 24
}

func (m *model) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune

	if y < 0 || y >= len(m.stats) {
		return ch, StyleNormal, nil, 1
	}

	if x >= 0 && x < 40 {
		if x < len(m.stats[y].Name) {
			ch = rune(m.stats[y].Name[x])
		} else {
			ch = ' '
		}
	} else if x >= 40 && x < 42 {
		ch = ' '
	} else {
		relx := x - 42
		sum := m.summaries[m.stats[y].ID]
		if relx < len(sum) {
			ch = rune(sum[relx])
		} else {
			ch = ' '
		}
	}
	/*
		style = m.styles[y]
		if m.items[y] == m.selected {
			style = style.Reverse(true)
		}
	*/
	// return ch, style, nil, 1
	return ch, StyleNormal, nil, 1
}
func (m *model) SetCursor(int, int)                {}
func (m *model) GetCursor() (int, int, bool, bool) { return 0, 0, false, false }
func (m *model) MoveCursor(offx, offy int)         {}

type statSummary struct {
	stat    intr.Stat
	summary string
}

type background struct {
	app            *views.Application
	title          *views.TextBar
	main           *views.CellView
	clockUps       chan string
	stats          chan []intr.Stat
	refreshSummary chan intr.Stat
	model          *model
	datasets       map[string]intr.Dataset
	sync.Mutex
}

func newBackground(app *views.Application, title *views.TextBar, m *views.CellView) *background {
	b := &background{
		app:            app,
		title:          title,
		clockUps:       make(chan string, 1),
		stats:          make(chan []intr.Stat, 1),
		refreshSummary: make(chan intr.Stat, 100),
		main:           m,
		model:          newModel(),
		datasets:       make(map[string]intr.Dataset),
	}
	b.main.SetModel(b.model)

	go b.loop()
	go b.clockTick()
	go b.list()
	for i := 0; i < 10; i++ {
		go b.summary()
	}

	return b
}

func (b *background) loop() {
	refresh := false
	for {
		refresh = false
		select {
		case t := <-b.clockUps:
			b.title.SetRight(t, tcell.StyleDefault)
			refresh = true
		case s := <-b.stats:
			b.model.SetStats(s)
			for _, x := range s[:100] {
				b.refreshSummary <- x
			}
			refresh = true
		}

		if refresh {
			b.app.Refresh()
		}
	}
}

func (b *background) clockTick() {
	for {
		b.clockUps <- time.Now().Format("15:04:05")
		time.Sleep(1 * time.Second)
	}
}

func (b *background) list() {
	for {
		stats, err := intr.StatList()
		if err != nil {
			fmt.Printf("StatList err: %s", err)
		}
		sort.Sort(intr.ByDataReceivedAt(stats))

		b.stats <- stats
		time.Sleep(time.Minute)
	}
}

func (b *background) summary() {
	for s := range b.refreshSummary {
		dset, err := intr.LoadDataset(s.ID, "1d15m1w")
		if err != nil {
			log.Printf("load dataset err: %s", err)
		}
		b.Lock()
		b.datasets[s.ID] = dset
		b.Unlock()
		sum := NewSummary(&dset)
		sumStr := fmt.Sprintf("%g\t%g\t%g\t%g\t%g\t%g\t%g\t%g", sum.Latest, sum.Min, sum.Max, sum.Mean, sum.Total, sum.StdDev, sum.Conf95Min, sum.Conf99Max)
		b.model.SetSummary(s.ID, sumStr)
	}
}

func top(cmd *cobra.Command, args []string) error {
	log.Printf("top config.Host: %s", config.Host())
	app := &views.Application{}
	window := &mainWindow{parent: app}

	title := views.NewTextBar()
	title.SetStyle(tcell.StyleDefault.Background(tcell.ColorTeal).Foreground(tcell.ColorWhite))
	title.SetCenter("StatHat top", tcell.StyleDefault)
	title.SetRight("21:22:37", tcell.StyleDefault)

	window.main = views.NewCellView()
	window.main.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))

	window.SetTitle(title)
	window.SetContent(window.main)

	app.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack))
	app.SetRootWidget(window)

	b := newBackground(app, title, window.main)
	_ = b

	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	return nil
}
