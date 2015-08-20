package main

import (
    "github.com/nsf/termbox-go"
    //"encoding/json"
    //"net/http"
    //"io/ioutil"
    //"strings"
)

type SeriesTab struct {
    EntriesTab
}

func newSeriesTab(a *Apollo) *SeriesTab {
    t := &SeriesTab{
        EntriesTab: EntriesTab{
            a: a,
            entries: &a.d.Series,
            name: "series",
            sortField: "title",
            view: "passive",
            entryType: "episodic",
        },
    }

    t.refreshSlice()

    return t
}

func (t *SeriesTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case 't':
        return true
    }

    return t.handleKeyEvent(ev)
}


func (t *SeriesTab) Query(query string) {
    t.query(query)

    if query[0] != ':' && t.a.c.get("autotag") == "true" {
        t.a.inputActive = false
    }
}
