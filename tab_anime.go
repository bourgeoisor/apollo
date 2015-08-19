package main

import (
    "github.com/nsf/termbox-go"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "strings"
    "strconv"
)

type AnimeTab struct {
    EntriesTab
}

func newAnimeTab(a *Apollo) *AnimeTab {
    t := &AnimeTab{
        EntriesTab: EntriesTab{
            a: a,
            entries: &a.d.Anime,
            name: "anime",
            sortField: "title",
            status: "anime",
            view: "passive",
        },
    }

    t.refreshSlice()

    return t
}

func (t *AnimeTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case 't':
        t.fetchHummingbirdTags()
        return true
    }

    return t.handleKeyEvent(ev)
}


func (t *AnimeTab) Query(query string) {
    t.query(query)

    if query[0] != ':' && t.a.c.get("autotag") == "true" {
        t.fetchHummingbirdTags()
        t.a.inputActive = false
    }
}

type HummingbirdEntry struct {
    Id int
    Title string
    Episode_count int
    Started_airing string
}

func (t *AnimeTab) fetchHummingbirdTags() {
    title := strings.Replace(t.slice[t.cursor].Title, " ", "+", -1)
    url := "http://hummingbird.me/api/v1/search/anime?query=" + title
    t.a.logDebug(url)

    res, err := http.Get(url)
    if err != nil {
        t.a.logError(err.Error())
        return
    }
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    var data []HummingbirdEntry
    err = json.Unmarshal(body, &data)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    for i := 0; i < len(data); i++ {
        if i < 10 {
            releaseDate := strings.Split(data[i].Started_airing, "-")
            t.search = append(t.search, Entry{
                Title: data[i].Title,
                TagID: strconv.Itoa(data[i].Id),
                Year: releaseDate[0],
                EpisodeTotal: data[i].Episode_count,
            })
        }
    }

    if len(t.search) > 0 {
        t.view = "tag"
    }
}
