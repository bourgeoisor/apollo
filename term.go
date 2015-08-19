package main

import (
    "github.com/nsf/termbox-go"
    "os"
    "log"
    "strconv"
    "unicode"
    "errors"
)

var colors map[rune]termbox.Attribute = map[rune]termbox.Attribute{
    'd': termbox.ColorDefault,
    'k': termbox.ColorBlack,    'K': termbox.ColorBlack | termbox.AttrBold,
    'r': termbox.ColorRed,      'R': termbox.ColorRed | termbox.AttrBold,
    'g': termbox.ColorGreen,    'G': termbox.ColorGreen | termbox.AttrBold,
    'y': termbox.ColorYellow,   'Y': termbox.ColorYellow | termbox.AttrBold,
    'b': termbox.ColorBlue,     'B': termbox.ColorBlue | termbox.AttrBold,
    'm': termbox.ColorMagenta,  'M': termbox.ColorMagenta | termbox.AttrBold,
    'c': termbox.ColorCyan,     'C': termbox.ColorCyan | termbox.AttrBold,
    'w': termbox.ColorWhite,    'W': termbox.ColorWhite | termbox.AttrBold,
}

type Tabber interface {
    Name() string
    Status() string
    HandleKeyEvent(*termbox.Event) bool
    Draw()
    Query(string)
}

type Apollo struct {
    running bool
    width int
    height int
    events chan termbox.Event
    c *Configuration
    d *Database
    currentTab int
    tabs []Tabber
    input []rune
    inputCursor int
    inputActive bool
}

func newApollo() *Apollo {
    err := os.Mkdir(os.Getenv("HOME") + "/.config/apollo", 0755)
    if err != nil {
        log.Print(err)
    }

    width, height := termbox.Size()
    var tabs []Tabber

    a := &Apollo{
        running: true,
        width: width,
        height: height,
        events: make(chan termbox.Event, 20),
        c: newConfiguration(),
        d: newDatabase(),
        tabs: tabs,
    }

    a.tabs = append(a.tabs, Tabber(newStatusTab(a)))

    if a.c.get("movies_tab") == "true" {
        a.tabs = append(a.tabs, Tabber(newMoviesTab(a)))
    }

    if a.c.get("series_tab") == "true" {
        a.tabs = append(a.tabs, Tabber(newSeriesTab(a)))
    }

    if a.c.get("anime_tab") == "true" {
        a.tabs = append(a.tabs, Tabber(newAnimeTab(a)))
    }

    if a.c.get("games_tab") == "true" {
        a.tabs = append(a.tabs, Tabber(newGamesTab(a)))
    }

    if a.c.get("books_tab") == "true" {
        a.tabs = append(a.tabs, Tabber(newBooksTab(a)))
    }

    a.printWelcome()

    return a
}

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

func (a *Apollo) handleKeyEvent(ev *termbox.Event) {
    if !a.inputActive && ev.Mod != termbox.ModAlt {
        handled := a.tabs[a.currentTab].HandleKeyEvent(ev)
        if handled {
            return
        }
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
    default:
        if ev.Mod == termbox.ModAlt {
            indexes := map[rune]int{'1': 1, '2': 2, '3': 3,
                                    '4': 4, '5': 5, '6': 6,
                                    '7': 7, '8': 8, '9': 9,}
            if i, exist := indexes[ev.Ch]; exist {
                if len(a.tabs) > i - 1 {
                    a.currentTab = i - 1
                }
            }
        } else {
            if unicode.IsPrint(ev.Ch) && a.inputActive {
                a.input = append(a.input, ' ')
                copy(a.input[a.inputCursor+1:], a.input[a.inputCursor:])
                a.input[a.inputCursor] = ev.Ch
                a.inputCursor++
            }
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
        }
    }
}

func (a *Apollo) drawStatusBars() {
    for i := 0; i < a.width; i++ {
        termbox.SetCell(i, 0, ' ', colors['d'], colors['k'])
        termbox.SetCell(i, a.height - 2, ' ', colors['d'], colors['k'])
    }
}

func (a *Apollo) drawTopStatus() {
    runes := []rune(version + " - " + a.tabs[a.currentTab].Status())
    for i := 0; i < len(runes); i++ {
        termbox.SetCell(i, 0, runes[i], colors['W'], colors['k'])
    }
}

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
        termbox.SetCell(x, a.height - 2, runes[i], fg, colors['k'])
        x++
    }
}

func (a *Apollo) drawInput() {
    if len(a.input) < a.width {
        for i := 0; i < len(a.input); i++ {
            termbox.SetCell(i, a.height - 1, a.input[i], colors['w'], colors['d'])
        }
    } else {
        offset := len(a.input) - a.width + 1
        for i := 0; i < a.width - 1; i++ {
            termbox.SetCell(i, a.height - 1, a.input[i + offset], colors['w'], colors['d'])
        }
    }

    if a.inputActive {
        termbox.SetCursor(a.inputCursor, a.height - 1)
    } else {
        termbox.HideCursor()
    }
}

func (a *Apollo) draw() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    a.tabs[a.currentTab].Draw()

    a.drawStatusBars()
    a.drawTopStatus()
    a.drawBottomStatus()
    a.drawInput()

    termbox.Flush()
}

func (a *Apollo) log(str string) {
    a.tabs[0].Query(str)
}

func (a *Apollo) logError(str string) {
    a.log("{r}│ ERROR: {d}" + str)
    log.Print(str)
}

func (a *Apollo) logDebug(str string) {
    if a.c.get("debug") == "true" {
        log.Print(str)
    }
}

func (a *Apollo) openTab(name string) error {
    for i := range a.tabs {
        if a.tabs[i].Name() == name {
            a.currentTab = i
            return nil
        }
    }

    switch name {
    case "movies":
        a.tabs = append(a.tabs, Tabber(newMoviesTab(a)))
    case "series":
        a.tabs = append(a.tabs, Tabber(newSeriesTab(a)))
    case "anime":
        a.tabs = append(a.tabs, Tabber(newAnimeTab(a)))
    case "games":
        a.tabs = append(a.tabs, Tabber(newGamesTab(a)))
    case "books":
        a.tabs = append(a.tabs, Tabber(newBooksTab(a)))
    default:
        return errors.New("term: tab does not exist")
    }

    a.currentTab = len(a.tabs) - 1
    return nil
}

func (a *Apollo) closeCurrentTab() error {
    if a.tabs[a.currentTab].Name() == "(status)" {
        return errors.New("term: cannot close status tab")
    }

    a.tabs = append(a.tabs[:a.currentTab], a.tabs[a.currentTab+1:]...)
    a.currentTab--
    return nil
}
