package main

import (
    "github.com/nsf/termbox-go"
    "log"
)

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
            a.draw()
            if err != nil {
                a.running = false
                log.Fatal(err)
            }
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
