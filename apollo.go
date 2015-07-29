package main

import (
    "github.com/nsf/termbox-go"
    "strconv"
    "unicode"
    "log"
)

type Tab interface {
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
    configuration *Configuration
    database *Database
    currentTab int
    tabs []Tab
    input []rune
    inputCursor int
    inputActive bool
}

func createApollo() *Apollo {
    width, height := termbox.Size()

    var tabs []Tab

    a := &Apollo{
        running: true,
        width: width,
        height: height,
        events: make(chan termbox.Event, 20),
        configuration: createConfiguration(),
        database: createDatabase(),
        tabs: tabs,
    }

    a.tabs = append(tabs, Tab(createStatusTab(a)))

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
    handled := false
    if !a.inputActive {
        handled = a.tabs[a.currentTab].HandleKeyEvent(ev)
    }

    if !handled {
        switch ev.Key {
        case termbox.KeyCtrlC:
            a.running = false
        case termbox.KeyEnter:
            if len(a.input) > 0 {
                if a.input[0] == '/' {
                    a.handleCommand()
                } else {
                    a.tabs[a.currentTab].Query(string(a.input))
                }
                a.input = a.input[:0]
                a.inputCursor = 0
            } else {
                if a.inputActive {
                    a.inputActive = false
                } else {
                    a.inputActive = true
                }
            }
        case termbox.KeyBackspace, termbox.KeyBackspace2:
            if a.inputActive {
                if a.inputCursor > 0 {
                    a.input = append(a.input[:a.inputCursor-1], a.input[a.inputCursor:]...)
                    a.inputCursor--
                }
            }
        case termbox.KeySpace:
            if a.inputActive {
                a.input = append(a.input, ' ')
                copy(a.input[a.inputCursor+1:], a.input[a.inputCursor:])
                a.input[a.inputCursor] = ' '
                a.inputCursor++
            }
        case termbox.KeyArrowLeft:
            if a.inputActive {
                a.inputCursor--
                if a.inputCursor < 0 {
                    a.inputCursor = 0
                }
            }
        case termbox.KeyArrowRight:
            if a.inputActive {
                a.inputCursor++
                if a.inputCursor > len(a.input) {
                    a.inputCursor = len(a.input)
                }
            }
        default:
            if ev.Mod == termbox.ModAlt {
                numbers := map[rune]int{'1': 1, '2': 2, '3': 3,
                                        '4': 4, '5': 5, '6': 6,
                                        '7': 7, '8': 8, '9': 9,}
                if number, exist := numbers[ev.Ch]; exist {
                    if len(a.tabs) > number - 1 {
                        a.currentTab = number - 1
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
    }
}

func (a *Apollo) draw() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    runes := []rune("top status bar - " + a.tabs[a.currentTab].Status())
    for i := 0; i < a.width; i++ {
        if i < len(runes) {
            termbox.SetCell(i, 0, runes[i], termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack | termbox.AttrBold)
        } else {
            termbox.SetCell(i, 0, ' ', termbox.ColorDefault, termbox.ColorBlack | termbox.AttrBold)
        }
    }

    a.tabs[a.currentTab].Draw()

    for i := 0; i < a.width; i++ {
        termbox.SetCell(i, a.height - 2, ' ', termbox.ColorDefault, termbox.ColorBlack | termbox.AttrBold)
    }
    x := 0
    for i := 0; i < len(a.tabs); i++ {
        runes := []rune(strconv.Itoa(i+1) + "." + a.tabs[i].Name() + " ")
        for j := 0; j < len(runes); j++ {
            termbox.SetCell(x, a.height - 2, runes[j], termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlack | termbox.AttrBold)
            x++
        }
    }

    if len(a.input) < a.width {
        for i := 0; i < len(a.input); i++ {
            termbox.SetCell(i, a.height - 1, a.input[i], termbox.ColorWhite, termbox.ColorDefault)
        }
    } else {
        offset := len(a.input) - a.width + 1
        for i := 0; i < a.width - 1; i++ {
            termbox.SetCell(i, a.height - 1, a.input[i + offset], termbox.ColorWhite, termbox.ColorDefault)
        }
    }
    if a.inputActive {
        termbox.SetCursor(a.inputCursor, a.height - 1)
    } else {
        termbox.HideCursor()
    }

    termbox.Flush()
}

func (a *Apollo) loop() {
    go func() {
        for {
            a.events <- termbox.PollEvent()
        }
    }()
    for (a.running) {
        select {
        case ev := <-a.events:
            err := a.handleEvent(&ev)
            if err != nil {
                a.running = false
                log.Fatal(err)
            }
            a.draw()
        }
    }
}

func main() {
    err := termbox.Init()
    if err != nil {
        log.Fatal(err)
    }
    defer termbox.Close()

    termbox.SetInputMode(termbox.InputAlt)

    apollo := createApollo()
    apollo.draw()
    apollo.loop()
}
