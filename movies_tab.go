package main

import (
    "github.com/nsf/termbox-go"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "strconv"
    "strings"
    "log"
    "sort"
)

type Entries struct {
    Title string
    Year string
    ImdbID string
}

type Data struct {
    Search []Entries
}

type MoviesTab struct {
    a *Apollo
    name string
    status string

    movies []Movie
    view string
    sorter string
    offset int
    cursor int
    omdb Data
    ratings bool
}

func CreateMoviesTab(a *Apollo) *MoviesTab {
    t := &MoviesTab{
        a: a,
        name: "movies",
        status: "movies",
        view: "watched",
        sorter: "title",
    }

    t.refreshSlice()

    return t
}

func (t *MoviesTab) Name() string {
    return t.name
}

func (t *MoviesTab) Status() string {
    return t.status
}

func (t *MoviesTab) HandleKeyEvent(ev *termbox.Event) bool {
    if len(t.omdb.Search) > 0 {
        indexes := map[rune]int{'0': 0, '1': 1, '2': 2, '3': 3,
                                '4': 4, '5': 5, '6': 6, '7': 7,
                                '8': 8, '9': 9,}
        if index, exist := indexes[ev.Ch]; exist {
            if index < len(t.omdb.Search) {
                log.Print(t.omdb.Search[index])
                dbIndex := t.index()
                t.a.d.Movies[dbIndex].Year = t.omdb.Search[index].Year
                t.a.d.Movies[dbIndex].Title = t.omdb.Search[index].Title
                t.a.d.Movies[dbIndex].ImdbID = t.omdb.Search[index].ImdbID
            }
        }

        t.omdb.Search = t.omdb.Search[:0]
        t.a.d.save()
        t.refreshSlice()

        return true
    }

    if t.view == "edit" {
        switch ev.Ch {
        case '0':
            t.a.inputActive = true
            t.a.input = []rune(":t " + t.a.d.Movies[t.index()].Title)
        case '1':
            t.a.inputActive = true
            t.a.input = []rune(":y " + t.a.d.Movies[t.index()].Year)
        }
        return true
    }

    switch ev.Ch {
    case '1':
        t.view = "watched"
        t.cursor = 0
        t.offset = 0
        t.refreshSlice()
    case '2':
        t.view = "unwatched"
        t.cursor = 0
        t.offset = 0
        t.refreshSlice()
    case '3':
        t.view = "all"
        t.cursor = 0
        t.offset = 0
        t.refreshSlice()
    case 's':
        if t.sorter == "title" {
            t.sorter = "year"
        } else if t.sorter == "year" {
            t.sorter = "rating"
        } else if t.sorter == "rating" {
            t.sorter = "title"
        }
        t.cursor = 0
        t.offset = 0
        t.sort()
    case 'D':
        if len(t.movies) > 0 {
            t.a.d.Movies = append(t.a.d.Movies[:t.index()], t.a.d.Movies[t.index()+1:]...)
            t.a.d.save()
            t.refreshSlice()
        }
    case 't':
        t.autoTag()
    case 'e':
        if t.view != "edit" {
            t.view = "edit"
        } else {
            t.view = "watched"
            t.refreshSlice()
        }
    case 'r':
        if t.ratings {
            t.ratings = false
        } else {
            t.ratings = true
        }
    case 'a':
        if t.a.d.Movies[t.index()].State == "Watched" {
            t.a.d.Movies[t.index()].State = "Unwatched"
        } else {
            t.a.d.Movies[t.index()].State = "Watched"
        }
        t.a.d.Movies[t.index()].Rating = 0
        t.a.d.save()
        t.refreshSlice()
    case 'z':
        if len(t.movies) > 0 {
            if t.a.d.Movies[t.index()].Rating > 0 {
                t.a.d.Movies[t.index()].Rating--
                t.a.d.save()
                t.refreshSlice()
            }
        }
    case 'x':
        if len(t.movies) > 0 {
            if t.a.d.Movies[t.index()].Rating < 6 {
                t.a.d.Movies[t.index()].Rating++
                t.a.d.save()
                t.refreshSlice()
            }
        }
    }

    switch ev.Key {
    case termbox.KeyArrowUp:
        t.cursor--
        if t.cursor < 0 {
            t.cursor = 0
        } else if t.cursor - t.offset < 0 {
            t.offset--
        }
    case termbox.KeyPgup:
        t.cursor -= 5
        if t.cursor < 0 {
            t.cursor = 0
            t.offset = 0
        } else if t.cursor - t.offset < 0 {
            t.offset -= 5
            if t.offset < 0 {
                t.offset = 0
            }
        }
    case termbox.KeyArrowDown:
        t.cursor++
        if t.cursor > len(t.movies) - 1 {
            t.cursor--
        } else if t.cursor - t.offset > t.a.height - 4 {
            t.offset++
        }
    case termbox.KeyPgdn:
        t.cursor += 5
        if t.cursor > len(t.movies) - 1 {
            t.cursor = len(t.movies) - 1
            if len(t.movies) > t.a.height - 3 {
                t.offset = len(t.movies) - (t.a.height - 3)
            }
        } else if t.cursor - t.offset > t.a.height - 4 {
            t.offset += 5
            if t.cursor - t.offset < t.a.height - 3 {
                t.offset = t.cursor - (t.a.height - 4)
            }
        }
    default:
        return false
    }

    log.Print("cursor: " + strconv.Itoa(t.cursor) + " offset: " + strconv.Itoa(t.offset) + 
    " len:" + strconv.Itoa(len(t.movies)) + " height:" + strconv.Itoa(t.a.height))

    return true
}

func (t *MoviesTab) Draw() {
    if t.view != "edit" {
        if len(t.omdb.Search) == 0 {
            for j := 0; j < t.a.height - 3; j++ {
                if j < len(t.movies) {
                    if t.ratings {
                        for i := 0; i < t.movies[j + t.offset].Rating; i++ {
                            if t.movies[j + t.offset].State == "Watched" {
                                termbox.SetCell(i + 3, j + 1, '*', color['y'], color['d'])
                            } else {
                                termbox.SetCell(i + 3, j + 1, '*', color['B'], color['d'])
                            }
                        }
                    }

                    runes := []rune(t.movies[j + t.offset].Year + " " + t.movies[j + t.offset].Title)
                    for i := 0; i < len(runes); i++ {
                        fg := color['d']
                        if i < 4 {
                            if t.movies[j + t.offset].State == "Watched" {
                                fg = color['g']
                            } else {
                                fg = color['b']
                            }
                        }

                        if t.ratings {
                            termbox.SetCell(i + 10, j + 1, runes[i], fg, color['d'])
                        } else {
                            termbox.SetCell(i + 3, j + 1, runes[i], fg, color['d'])
                        }
                    }
                }
            }

            termbox.SetCell(1, t.cursor - t.offset + 1, '*', color['d'], color['d'])
        } else {
            for j := 0; j < len(t.omdb.Search); j++ {
                runes := []rune(strconv.Itoa(j) + ". [" + t.omdb.Search[j].Year + "] " + t.omdb.Search[j].Title)
                for i := 0; i < len(runes); i++ {
                    termbox.SetCell(i, j + 1, runes[i], color['d'], color['d'])
                }
            }
        }
    } else {
        runes := []rune("0. " + t.a.d.Movies[t.index()].Title)
        for i := 0; i < len(runes); i++ {
            termbox.SetCell(i, 1, runes[i], color['d'], color['d'])
        }

        runes = []rune("1. " + t.a.d.Movies[t.index()].Year)
        for i := 0; i < len(runes); i++ {
            termbox.SetCell(i, 2, runes[i], color['d'], color['d'])
        }
    }
}

func (t *MoviesTab) Query(query string) {
    if query[0] != ':' {
        t.a.d.Movies = append(t.a.d.Movies, Movie{Title: query, State: "Watched"})
        t.a.d.save()
        t.refreshSlice()

        for i := 0; i < len(t.movies); i++ {
            if t.movies[i].Title == query {
                t.cursor = i
                if t.cursor > t.a.height - 4 {
                    t.offset = t.cursor - (t.a.height - 4)
                }
                t.a.inputActive = false
                t.autoTag()
            }
        }
    } else {
        if query[1] == 't' {
            t.a.d.Movies[t.index()].Title = query[3:]
        }

        if query[1] == 'y' {
            t.a.d.Movies[t.index()].Year = query[3:]
        }
    }
}

func (t *MoviesTab) index() int {
    for i := 0; i < len(t.a.d.Movies); i++ {
        if t.a.d.Movies[i].Title == t.movies[t.cursor].Title {
            return i
        }
    }

    return -1
}

func (t *MoviesTab) refreshSlice() {
    t.movies = t.movies[:0]
    for i := 0; i < len(t.a.d.Movies); i++ {
        if t.a.d.Movies[i].State == "Watched" {
            if t.view == "watched" || t.view == "all" {
                t.movies = append(t.movies, t.a.d.Movies[i])
            }
        } else if t.a.d.Movies[i].State == "Unwatched" {
            if t.view == "unwatched" || t.view == "all" {
                t.movies = append(t.movies, t.a.d.Movies[i])
            }
        }
    }
    t.sort()

    if t.cursor > len(t.movies) - 1 {
        t.cursor--
    }

    if t.offset > 0 {
        if len(t.movies) - t.offset < t.a.height - 3 {
            t.offset--
        }
    }

    t.status = "movies - " + t.view + " (" + strconv.Itoa(len(t.movies)) + " entries)"
}

func (t *MoviesTab) sort() {
    var titles []string
    for i := 0; i < len(t.movies); i++ {
        titles = append(titles, t.movies[i].Title)
    }
    sort.Strings(titles)

    var movies []Movie
    for j := 0; j < len(titles); j++ {
        for i := 0; i < len(t.movies); i++ {
            if titles[j] == t.movies[i].Title {
                movies = append(movies, t.movies[i])
            }
        }
    }
    t.movies = movies

    if t.sorter == "year" {

    } else if t.sorter == "rating" {
        t.sortByRating()
    }
}

func (t *MoviesTab) sortByRating() {
    var movies []Movie
    for j := 0; j < 7; j++ {
        for i := 0; i < len(t.movies); i++ {
            if t.movies[i].Rating == 6 - j {
                movies = append(movies, t.movies[i])
            }
        }
    }
    t.movies = movies
}

func (t *MoviesTab) autoTag() {
    title := strings.Replace(t.movies[t.cursor].Title, " ", "+", -1)
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

    err = json.Unmarshal(body, &t.omdb)
    if err != nil {
        log.Fatal(err)
    } else {
        for i := 0; i < len(t.omdb.Search); i++ {
            log.Print(t.omdb.Search[i].Title + " - " + t.omdb.Search[i].Year)
        }
        log.Print(t.omdb)
    }
}
