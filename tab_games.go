package main

import (
    "github.com/nsf/termbox-go"
    "encoding/xml"
    "net/http"
    "io/ioutil"
    "strings"
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
            view: "passive",
            additionalField: "platform",
            entryType: "additional",
        },
    }

    t.refreshSlice()

    return t
}

func (t *GamesTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case 't':
        t.fetchGamesDBTags()
        return true
    }

    return t.handleKeyEvent(ev)
}


func (t *GamesTab) Query(query string) {
    t.query(query)

    if query[0] != ':' && t.a.c.get("autotag") == "true" {
        t.a.inputActive = false
        t.fetchGamesDBTags()
    }
}

type Game struct {
    id string
    GameTitle string
    ReleaseDate string
    Platform string
}

type GamesDBData struct {
    XMLName xml.Name `xml:"Data"`
    Game []Game
}

func (t *GamesTab) fetchGamesDBTags() {
    title := strings.Replace(t.slice[t.cursor].Title, " ", "+", -1)
    url := "http://thegamesdb.net/api/GetGamesList.php?name=" + title
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

    var data GamesDBData
    err = xml.Unmarshal(body, &data)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    platforms := map[string]string{
        "Sony Playstation": "PS", "Sony Playstation 2": "PS2",
        "Sony Playstation 3": "PS3", "Sony Playstation 4": "PS4",
        "Sony PSP": "PSP", "Sony Playstation Vita": "VITA",
        "Microsoft Xbox": "XBOX", "Microsoft Xbox 360": "X360",
        "Microsoft Xbox One": "XONE",
        "Nintendo Entertainment System (NES)": "NES",
        "Super Nintendo (SNES)": "SNES", "Nintendo 64": "N64",
        "Nintendo GameCube": "NGC", "Nintendo DS": "NDS", "Nintendo 3DS": "3DS",
        "Nintendo Game Boy": "GB",
        "Nintendo Game Boy Color": "GBC", "Nintendo Game Boy Advance": "GBA",
        "Nintendo Wii": "WII", "Nintendo Wii U": "WIIU",
        "PC": "PC",
    }

    for i := 0; i < len(data.Game); i++ {
        if i < 10 {
            releaseDate := strings.Split(data.Game[i].ReleaseDate, "/")
            t.search = append(t.search, Entry{
                Title: data.Game[i].GameTitle,
                Year: releaseDate[2],
                TagID: data.Game[i].id,
                Info1: platforms[data.Game[i].Platform],
            })
        }
    }

    if len(t.search) > 0 {
        t.view = "tag"
    }
}
