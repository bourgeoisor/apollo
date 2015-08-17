package main

import (
    "github.com/nsf/termbox-go"
)

type StatusTab struct {
    a *Apollo
    name string
    status string

    history []string
    offset int
}

func newStatusTab(a *Apollo) *StatusTab {
    t := &StatusTab{
        a: a,
        name: "(status)",
        status: "logs",
        history: make([]string, 200),
    }

    return t
}

func (t *StatusTab) Name() string {
    return t.name
}

func (t *StatusTab) Status() string {
    return t.status
}

func (t *StatusTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Key {
    case termbox.KeyPgup:
        t.offset += 5
        if t.offset > 200 - t.a.height + 3 {
            t.offset = 200 - t.a.height + 3
        }
    case termbox.KeyPgdn:
        t.offset -= 5
        if t.offset < 0 {
            t.offset = 0
        }
    default:
        return false
    }

    return true
}

func (t *StatusTab) Draw() {
    historySlice := t.history[200-t.a.height+3-t.offset:200-t.offset]

    for j := 1; j < t.a.height - 3; j++ {
        fg := colors['d']
        x := 0
        runes := []rune(historySlice[j])
        for i := 0; i < len(runes); i++ {
            if runes[i] == '{' {
                fg = colors[runes[i+1]]
                i += 3
            }
            termbox.SetCell(x, j, runes[i], fg, colors['d'])
            x++
        }
    }
}

func (t *StatusTab) Query(query string) {
    t.history = t.history[1:]
    t.history = append(t.history, query)
}
