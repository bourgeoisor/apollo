package main

import (
    "github.com/nsf/termbox-go"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "strings"
    "log"
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
            additionalField: "author",
        },
    }

    t.refreshSlice()

    return t
}

func (t *BooksTab) HandleKeyEvent(ev *termbox.Event) bool {
    switch ev.Ch {
    case 't':
        t.fetchGoogleBooksTags()
        return true
    }

    return t.handleKeyEvent(ev)
}


func (t *BooksTab) Query(query string) {
    t.query(query)

    if query[0] != ':' && t.a.c.get("autotag") == "true" {
        t.fetchGoogleBooksTags()
        t.a.inputActive = false
    }
}

type GoogleBooksInfo struct {
    Title string
    Authors []string
    PublishedDate string
}

type GoogleBooksEntry struct {
    Id string
    VolumeInfo GoogleBooksInfo
}

type GoogleBooksData struct {
    Items []GoogleBooksEntry
}

func (t *BooksTab) fetchGoogleBooksTags() {
    title := strings.Replace(t.slice[t.cursor].Title, " ", "+", -1)
    url := "https://www.googleapis.com/books/v1/volumes?q=" + title + "&projection=lite&printType=books&maxResults=10"
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

    var data GoogleBooksData
    err = json.Unmarshal(body, &data)
    if err != nil {
        t.a.logError(err.Error())
        return
    }

    log.Print(data)

    for i := 0; i < len(data.Items); i++ {
        if i < 10 {
            if len(data.Items[i].VolumeInfo.Authors) == 0 {
                data.Items[i].VolumeInfo.Authors = append(data.Items[i].VolumeInfo.Authors, "")
            }
            t.search = append(t.search, Entry{
                Title: data.Items[i].VolumeInfo.Title,
                TagID: data.Items[i].Id,
                Info1: data.Items[i].VolumeInfo.Authors[0],
            })
        }
    }

    if len(t.search) > 0 {
        t.view = "tag"
    }
}
