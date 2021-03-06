package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/spf13/cobra"
	"github.com/stathat/cmd/stathat/config"
	"github.com/stathat/cmd/stathat/intr"
	"github.com/stathat/numbers"
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
	StyleBold   = tcell.StyleDefault.Foreground(tcell.ColorSilver).Background(tcell.ColorBlack).Bold(true)
	StyleGood   = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	StyleWarn   = tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack)
	StyleError  = tcell.StyleDefault.Foreground(tcell.ColorMaroon).Background(tcell.ColorBlack)
)

type mainWindow struct {
	parent     *views.Application
	view       views.View
	main       *views.CellView
	keybar     *views.SimpleStyledText
	status     *views.SimpleStyledTextBar
	background *background

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
			case 'R', 'r':
				a.background.SortRecent()
				return true
			case 'D', 'd':
				a.background.SortData()
				return true
			case ' ':
				a.background.Refresh()
				return true
			}

		}
	}

	return a.Panel.HandleEvent(ev)
}

type model struct {
	stats     []intr.Stat
	summaries map[string]string
	timeframe string
	heading   string
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

func (m *model) SetTimeframe(t string) {
	m.Lock()
	m.timeframe = t
	m.heading = fmt.Sprintf("%-41s%9s%9s%9s%9s%9s%9s%9s%9s", m.timeframe, "Latest", "Min", "Max", "Mean", "Total", "StdDev", "95% Min", "95% Max")
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

	if y == 0 {
		if x >= 0 && x < len(m.heading) {
			ch = rune(m.heading[x])
		} else {
			ch = ' '
		}

		return ch, StyleBold, nil, 1
	}

	rely := y - 1
	if x >= 0 && x < 40 {
		if x < len(m.stats[rely].Name) {
			ch = rune(m.stats[rely].Name[x])
		} else {
			ch = ' '
		}
	} else if x >= 40 && x < 42 {
		ch = ' '
	} else {
		relx := x - 42
		sum := m.summaries[m.stats[rely].ID]
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

type sortmeth int

const (
	sortData sortmeth = iota
	sortRecent
)

type background struct {
	app            *views.Application
	title          *views.TextBar
	main           *views.CellView
	clockUps       chan string
	stats          chan []intr.Stat
	refreshSummary chan intr.Stat
	model          *model
	datasets       map[string]intr.Dataset
	timeframe      string
	refresh        chan bool
	sortMethod     sortmeth
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
		timeframe:      "1d15m1w",
		refresh:        make(chan bool),
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

func (b *background) Refresh() {
	b.refresh <- true
}

func (b *background) SortRecent() {
	b.Lock()
	b.sortMethod = sortRecent
	b.Unlock()
	b.Refresh()
}

func (b *background) SortData() {
	b.Lock()
	b.sortMethod = sortData
	b.Unlock()
	b.Refresh()
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
			b.model.SetTimeframe(b.timeframe)
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
		switch b.sortMethod {
		case sortData:
			sort.Sort(intr.ByDataReceivedAt(stats))
		case sortRecent:
			sort.Sort(intr.ByCreatedAt(stats))
		default:
			sort.Sort(intr.ByDataReceivedAt(stats))
		}

		b.stats <- stats
		select {
		case <-b.refresh:
			continue
		case <-time.After(time.Minute):
		}
	}
}

func (b *background) summary() {
	for s := range b.refreshSummary {
		dset, err := intr.LoadDataset(s.ID, b.timeframe)
		if err != nil {
			log.Printf("load dataset err: %s", err)
		}
		b.Lock()
		b.datasets[s.ID] = dset
		b.Unlock()
		width := 8
		sum := NewSummary(&dset)
		sumStrs := []string{
			floatCol(sum.Latest, width),
			floatCol(sum.Min, width),
			floatCol(sum.Max, width),
			floatCol(sum.Mean, width),
			floatCol(sum.Total, width),
			floatCol(sum.StdDev, width),
			floatCol(sum.Conf95Min, width),
			floatCol(sum.Conf95Max, width),
		}
		if !s.Counter {
			sumStrs[4] = fmt.Sprintf("%8s", "---")
		}

		sumStr := strings.Join(sumStrs, " ")
		b.model.SetSummary(s.ID, sumStr)
	}
}

func floatCol(x float64, width int) string {
	s := numbers.Humanize(x)

	return fmt.Sprintf("%8s", s)
}

// floatCol formats x to be 12 chars long.
// http://stackoverflow.com/questions/36515818/golang-is-there-any-standard-library-to-convert-float64-to-string-with-fix-widt
func floatCol2(x float64, width int) string {
	if x >= 1e12 {
		// Check to see how many fraction digits fit in:
		s := fmt.Sprintf("%.g", x)
		format := fmt.Sprintf("%%12.%dg", width-len(s))
		return fmt.Sprintf(format, x)
	}

	// Check to see how many fraction digits fit in:
	s := fmt.Sprintf("%.0f", x)
	if len(s) == width {
		return s
	}
	format := fmt.Sprintf("%%%d.%df", len(s), width-len(s)-1)
	return fmt.Sprintf(format, x)
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
	window.background = b

	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	return nil
}
