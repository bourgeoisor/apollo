package main

import (
    "github.com/nsf/termbox-go"
    //"encoding/json"
    //"net/http"
    //"io/ioutil"
    //"strings"
)

type BooksTab struct {
    EntriesTab
}

func newBooksTab(a *Apollo) *BooksTab {
    t := &BooksTab{
        EntriesTab: EntriesTab{
            a: a,
            entries: &a.d.Books,
            name: "books",
            sortField: "title",
            status: "books",
            view: "passive",
        },
    }

    t.refreshSlice()

    return t
}

func (t *BooksTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case 't':
        return true
    }

    return t.handleKeyEvent(ev)
}


func (t *BooksTab) Query(query string) {
    t.query(query)

    if query[0] != ':' && t.a.c.get("autotag") == "true" {
        t.a.inputActive = false
    }
}
