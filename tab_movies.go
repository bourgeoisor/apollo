package main

import (
    "github.com/nsf/termbox-go"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "strings"
    "log"
)

type OMDBEntry struct {
    Title string
    Year string
    ImdbID string
}

type OMDBData struct {
    Search []OMDBEntry
}

type MoviesTab struct {
    EntriesTab
    omdb OMDBData
}

func newMoviesTab(a *Apollo) *MoviesTab {
    t := &MoviesTab{
        EntriesTab: EntriesTab{
            a: a,
            entries: &a.d.Movies,
            name: "movies",
            sortField: "title",
            status: "movies",
            view: "passive",
        },
    }

    t.refreshSlice()

    return t
}

func (t *MoviesTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case 't':
        t.fetchOMDBTags()
        return true
    }

    return t.handleKeyEvent(ev)
}


func (t *MoviesTab) Query(query string) {
    t.query(query)

    if query[0] != ':' {
        t.a.inputActive = false
        t.fetchOMDBTags()
    }
}

func (t *MoviesTab) fetchOMDBTags() {
    title := strings.Replace(t.slice[t.cursor].Title, " ", "+", -1)
    url := "http://www.omdbapi.com/?s=" + title + "&type=movie&y=&plot=full&r=json"
    log.Print(url)

    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    var data OMDBData
    err = json.Unmarshal(body, &data)
    if err != nil {
        log.Fatal(err)
    } else {
        for i := 0; i < len(data.Search); i++ {
            t.search = append(t.search, Entry{
                Title: data.Search[i].Title,
                Year: data.Search[i].Year,
                TagID: data.Search[i].ImdbID,
            })
        }
        log.Print(t.search)
    }

    if len(t.search) > 0 {
        t.view = "tag"
    }
}
