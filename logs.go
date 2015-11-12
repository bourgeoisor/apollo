package main

import (
	"github.com/nsf/termbox-go"
)

// StatusTab is a tab for displaying logs implementing Tabber.
type StatusTab struct {
	a      *Apollo
	name   string
	status string

	history []string
	offset  int
}

// NewStatusTab creates a new StatusTab and returns it.
func newStatusTab(a *Apollo) *StatusTab {
	t := &StatusTab{
		a:       a,
		name:    "(status)",
		status:  "logs",
		history: make([]string, 200),
	}

	return t
}

// Name returns the name of the tab.
func (t *StatusTab) Name() string {
	return t.name
}

// Status returns the status of the tab.
func (t *StatusTab) Status() string {
	return t.status
}

// HandleKeyEvent handles the changes in logs offsets.
func (t *StatusTab) HandleKeyEvent(ev *termbox.Event) {
	switch ev.Key {
	case termbox.KeyPgup:
		t.offset += 5
		if t.offset > 200-t.a.height+3 {
			t.offset = 200 - t.a.height + 3
		}
	case termbox.KeyPgdn:
		t.offset -= 5
		if t.offset < 0 {
			t.offset = 0
		}
	}
}

// Draw creates a slice of the logs history and draws it on the screen.
func (t *StatusTab) Draw() {
	historySlice := t.history[200-t.a.height+2-t.offset : 200-t.offset]

	for j := 1; j < t.a.height-2; j++ {
		t.a.drawString(0, j, historySlice[j])
	}
}

// Query adds a new string to the list of logs.
func (t *StatusTab) Query(query string) {
	if query[0] != '!' {
		t.history = t.history[1:]
		t.history = append(t.history, query)
	}
}
