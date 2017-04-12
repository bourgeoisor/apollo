package main

import (
	"errors"
	"github.com/nsf/termbox-go"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// Colors is a map of all the colors available through termbox.
var colors map[rune]termbox.Attribute = map[rune]termbox.Attribute{
	'd': termbox.ColorDefault,
	'k': termbox.ColorBlack, 'K': termbox.ColorBlack | termbox.AttrBold,
	'r': termbox.ColorRed, 'R': termbox.ColorRed | termbox.AttrBold,
	'g': termbox.ColorGreen, 'G': termbox.ColorGreen | termbox.AttrBold,
	'y': termbox.ColorYellow, 'Y': termbox.ColorYellow | termbox.AttrBold,
	'b': termbox.ColorBlue, 'B': termbox.ColorBlue | termbox.AttrBold,
	'm': termbox.ColorMagenta, 'M': termbox.ColorMagenta | termbox.AttrBold,
	'c': termbox.ColorCyan, 'C': termbox.ColorCyan | termbox.AttrBold,
	'w': termbox.ColorWhite, 'W': termbox.ColorWhite | termbox.AttrBold,
}

// Tabber is an interface used by the different tabs.
type Tabber interface {
	Name() string
	Status() string
	HandleKeyEvent(*termbox.Event)
	Draw()
	Query(string)
}

// Apollo is the main object of the application.
type Apollo struct {
	running     bool
	width       int
	height      int
	events      chan termbox.Event
	c           *Configuration
	d           *Database
	currentTab  int
	tabs        []Tabber
	input       []rune
	inputCursor int
	inputActive bool
}

// NewApollo creates a new Apollo, initializing a new Configuration and new Database in the process.
// It opens the default tabs and then returns itself.
func newApollo() *Apollo {
	width, height := termbox.Size()
	var tabs []Tabber

	a := &Apollo{
		running: true,
		width:   width,
		height:  height,
		events:  make(chan termbox.Event, 20),
		c:       newConfiguration(),
		d:       newDatabase(),
		tabs:    tabs,
	}

	a.tabs = append(a.tabs, Tabber(newStatusTab(a)))

	autoOpenTabs := strings.Split(a.c.get("tabs-startup"), ",")
	for _, name := range autoOpenTabs {
		a.openTab(name)
	}

	a.printWelcome()
	a.currentTab = 0

	return a
}

// HandleEvent changes the size of the terminal if it's been resized. On all other types of events,
// it sends them to the key handler. It returns in case of an error.
func (a *Apollo) handleEvent(ev *termbox.Event) error {
	switch ev.Type {
	case termbox.EventKey:
		a.handleKeyEvent(ev)
	case termbox.EventResize:
		a.width, a.height = termbox.Size()
	case termbox.EventError:
		return ev.Err
	}

	return nil
}

// HandleKeyEvent handles the current key event. It will handle events related to the input bar
// and the changes of current tab. If none of those happens, it'll forward the event to the
// current tab's event handler.
func (a *Apollo) handleKeyEvent(ev *termbox.Event) {
	if ev.Mod == termbox.ModAlt {
		indexes := map[rune]int{'1': 1, '2': 2, '3': 3,
			'4': 4, '5': 5, '6': 6,
			'7': 7, '8': 8, '9': 9}
		if i, exist := indexes[ev.Ch]; exist {
			if len(a.tabs) > i-1 {
				a.currentTab = i - 1
				a.tabs[a.currentTab].Query("!focused")
			}
		}
		return
	}

	switch ev.Key {
	case termbox.KeyCtrlC:
		a.running = false
	case termbox.KeyEnter:
		if len(a.input) > 0 {
			if a.input[0] == '/' {
				a.handleCommand()
			} else if a.currentTab != 0 {
				a.tabs[a.currentTab].Query(string(a.input))
			}
			a.input = a.input[:0]
			a.inputCursor = 0
		} else {
			a.inputActive = !a.inputActive
		}
	}

	if a.inputActive {
		switch ev.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if a.inputCursor > 0 {
				a.input = append(a.input[:a.inputCursor-1], a.input[a.inputCursor:]...)
				a.inputCursor--
			}
		case termbox.KeySpace:
			a.input = append(a.input, ' ')
			copy(a.input[a.inputCursor+1:], a.input[a.inputCursor:])
			a.input[a.inputCursor] = ' '
			a.inputCursor++
		case termbox.KeyArrowLeft:
			a.inputCursor--
			if a.inputCursor < 0 {
				a.inputCursor = 0
			}
		case termbox.KeyArrowRight:
			a.inputCursor++
			if a.inputCursor > len(a.input) {
				a.inputCursor = len(a.input)
			}
		default:
			if unicode.IsPrint(ev.Ch) {
				a.input = append(a.input, ' ')
				copy(a.input[a.inputCursor+1:], a.input[a.inputCursor:])
				a.input[a.inputCursor] = ev.Ch
				a.inputCursor++
			}
		}
	} else {
		a.tabs[a.currentTab].HandleKeyEvent(ev)
	}
}

// DrawString draws a given string on a given row.
func (a *Apollo) drawString(x, y int, str string) {
	fg := colors['d']
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '{' {
			fg = colors[runes[i+1]]
			i += 3
		}
		termbox.SetCell(x, y, runes[i], fg, colors['d'])
		x++
	}
}

// DrawStatusBars draws the background color of the two status rows.
func (a *Apollo) drawStatusBars() {
	for i := 0; i < a.width; i++ {
		termbox.SetCell(i, 0, ' ', colors['d'], colors['b'])
		termbox.SetCell(i, a.height-2, ' ', colors['d'], colors['b'])
	}
}

// DrawTopStatus draws the top status row.
func (a *Apollo) drawTopStatus() {
	runes := []rune(version + " - " + a.tabs[a.currentTab].Status())
	for i := 0; i < len(runes); i++ {
		termbox.SetCell(i, 0, runes[i], colors['W'], colors['b'])
	}
}

// DrawBottomStatus draws the tab status row of the terminal.
func (a *Apollo) drawBottomStatus() {
	var str string
	for i := range a.tabs {
		if i == a.currentTab {
			str += "{" + strconv.Itoa(i+1) + "." + a.tabs[i].Name() + "} "
		} else {
			str += strconv.Itoa(i+1) + "." + a.tabs[i].Name() + " "
		}
	}

	fg := colors['w']
	x := 0
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '{' {
			fg = colors['W']
			i++
		} else if runes[i] == '}' {
			fg = colors['w']
			i++
		}
		termbox.SetCell(x, a.height-2, runes[i], fg, colors['b'])
		x++
	}
}

// DrawInput draws the input row of the terminal.
func (a *Apollo) drawInput() {
	if len(a.input) < a.width {
		for i := 0; i < len(a.input); i++ {
			termbox.SetCell(i, a.height-1, a.input[i], colors['w'], colors['d'])
		}
	} else {
		offset := len(a.input) - a.width + 1
		for i := 0; i < a.width-1; i++ {
			termbox.SetCell(i, a.height-1, a.input[i+offset], colors['w'], colors['d'])
		}
	}

	if a.inputActive {
		termbox.SetCursor(a.inputCursor, a.height-1)
	} else {
		termbox.HideCursor()
	}
}

// Draw calls the different drawing functions for the terminal.
func (a *Apollo) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	a.tabs[a.currentTab].Draw()

	a.drawStatusBars()
	a.drawTopStatus()
	a.drawBottomStatus()
	a.drawInput()

	termbox.Flush()
}

// Log prints a message to the logs.
func (a *Apollo) log(str string) {
	a.tabs[0].Query(str)
}

// LogError prints a message to the logs and stderr.
func (a *Apollo) logError(str string) {
	a.log("{r}│ ERROR: {d}" + str)
	//log.Print(str)
}

// LogDebug logs the given string if the debug flag is on.
func (a *Apollo) logDebug(str string) {
	if a.c.get("debug") == "true" {
		log.Print(str)
	}
}

// OpenTab opens the given tab, if it's not already opened. If it is, it'll switch to it.
func (a *Apollo) openTab(name string) error {
	for i := range a.tabs {
		if a.tabs[i].Name() == name {
			a.currentTab = i
			return nil
		}
	}

	switch name {
	case "movies":
		a.tabs = append(a.tabs, Tabber(newEntriesTab(a, &a.d.Movies, "movies", "default", "", "omdb")))
	case "series":
		a.tabs = append(a.tabs, Tabber(newEntriesTab(a, &a.d.Series, "series", "episodic", "", "omdb")))
	case "anime":
		a.tabs = append(a.tabs, Tabber(newEntriesTab(a, &a.d.Anime, "anime", "episodic", "", "hummingbird")))
	case "games":
		a.tabs = append(a.tabs, Tabber(newEntriesTab(a, &a.d.Games, "games", "additional", "platform", "gamesdb")))
	case "books":
		a.tabs = append(a.tabs, Tabber(newEntriesTab(a, &a.d.Books, "books", "additional", "author", "googlebooks")))
	default:
		return errors.New("term: tab does not exist")
	}

	a.currentTab = len(a.tabs) - 1
	return nil
}

// CloseCurrentTab closes the tab currently being viewed, with the exception of the logs tab.
func (a *Apollo) closeCurrentTab() error {
	if a.tabs[a.currentTab].Name() == "(status)" {
		return errors.New("term: cannot close status tab")
	}

	a.tabs = append(a.tabs[:a.currentTab], a.tabs[a.currentTab+1:]...)
	a.currentTab--
	return nil
}
