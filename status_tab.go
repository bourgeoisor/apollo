package main

import (
    "github.com/nsf/termbox-go"
)

type StatusTab struct {
    a *Apollo
    name string

    history []string
}

func createStatusTab(a *Apollo) *StatusTab {
    t := &StatusTab{
        a: a,
        name: "status",
        history: make([]string, 200),
    }

    return t
}

func (t *StatusTab) Name() string {
    return t.name
}

func (t *StatusTab) Status() string {
    return "wat"
}

func (t *StatusTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Key {
    default:
        return false
    }

    return true
}

func (t *StatusTab) Draw() {
    historySlice := t.history[200-t.a.height+3:200]

    for j := 1; j < t.a.height - 3; j++ {
        runes := []rune(historySlice[j])
        for i := 0; i < len(runes); i++ {
            termbox.SetCell(i, j, runes[i], termbox.ColorWhite | termbox.AttrBold, termbox.ColorDefault)
        }
    }
}

func (t *StatusTab) Query(query string) {
    t.history = t.history[1:]
    t.history = append(t.history, query)
}
