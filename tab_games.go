package main

import (
    "github.com/nsf/termbox-go"
    //"encoding/json"
    //"net/http"
    //"io/ioutil"
    //"strings"
)

type GamesTab struct {
    EntriesTab
}

func newGamesTab(a *Apollo) *GamesTab {
    t := &GamesTab{
        EntriesTab: EntriesTab{
            a: a,
            entries: &a.d.Games,
            name: "games",
            sortField: "title",
            status: "games",
            view: "passive",
            additionalField: "console",
        },
    }

    t.refreshSlice()

    return t
}

func (t *GamesTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case 't':
        return true
    }

    return t.handleKeyEvent(ev)
}


func (t *GamesTab) Query(query string) {
    t.query(query)

    if query[0] != ':' && t.a.c.get("autotag") == "true" {
        t.a.inputActive = false
    }
}
