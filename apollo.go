package main

import (
    "github.com/nsf/termbox-go"
    "log"
)

type apollo struct {
    running bool
    width int
    height int
    events chan termbox.Event
    settings map[string]string
}

func create() *apollo {
    settings := map[string]string{
        "sett": "value",
    }

    width, height := termbox.Size()

    a := &apollo{
        running: true,
        width: width,
        height: height,
        events: make(chan termbox.Event, 20),
        settings: settings,
    }

    return a
}

func (a *apollo) handleEvent(ev *termbox.Event) error {
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

func (a *apollo) handleKeyEvent(ev *termbox.Event) {
    if ev.Key == termbox.KeyCtrlC {
        a.running = false
    }
}

func (a *apollo) draw() {

}

func (a *apollo) loop() {
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

    apollo := create()
    apollo.loop()
}
