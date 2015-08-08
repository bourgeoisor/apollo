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
    imdbID string
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
                t.a.d.Movies[t.index()].Year = t.omdb.Search[index].Year
                t.a.d.Movies[t.index()].Title = t.omdb.Search[index].Title
            }
        }

        t.omdb.Search = t.omdb.Search[:0]
        t.a.d.save()
        t.refreshSlice()

        return true
    }

    switch ev.Ch {
    case '1':
        t.view = "watched"
        t.refreshSlice()
    case '2':
        t.view = "unwatched"
        t.refreshSlice()
    case '3':
        t.view = "all"
        t.refreshSlice()
    case 's':
        t.sort()
    case 'D':
        if len(t.movies) > 0 {
            t.a.d.Movies = append(t.a.d.Movies[:t.index()], t.a.d.Movies[t.index()+1:]...)
            t.refreshSlice()
        }
    case 't':
        t.autoTag()
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

    log.Print("cursor: " + strconv.Itoa(t.cursor))

    return true
}

func (t *MoviesTab) Draw() {
    if len(t.omdb.Search) == 0 {
        for j := 0; j < t.a.height - 3; j++ {
            if j < len(t.movies) {
                runes := []rune("[" + t.movies[j + t.offset].Year + "] " + t.movies[j + t.offset].Title)
                for i := 0; i < len(runes); i++ {
                    termbox.SetCell(i + 3, j + 1, runes[i], color['d'], color['d'])
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
}

func (t *MoviesTab) Query(query string) {
    t.a.d.Movies = append(t.a.d.Movies, Movie{Title: query, State: "Watched"})
    t.a.d.save()
    t.refreshSlice()
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
    t.cursor = 0
    t.offset = 0
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
    }
}
